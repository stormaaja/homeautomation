package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"stormaaja/go-ha/data-store/configuration"
	"stormaaja/go-ha/data-store/configvalidators"
	"stormaaja/go-ha/data-store/dataroutes"
	"stormaaja/go-ha/data-store/genericroutes"
	v1 "stormaaja/go-ha/data-store/routes/v1"
	"stormaaja/go-ha/data-store/spot"

	"stormaaja/go-ha/data-store/store"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GetGinEnvironment() string {
	switch os.Getenv("ENVIRONMENT") {
	case "production":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}

func GetLogFile() string {
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		timeStamp := time.Now().Format("2006-01-02-15-04-05")
		logFile = fmt.Sprintf("data-store-log-%s.log", timeStamp)
	}
	return logFile
}

func CreateRoutes(
	memoryStore *store.MemoryStore,
	measurementStores []store.MeasurementStore,
	spotPriceApiClient *spot.SpotHintaApiClient,
	minerStateStore *store.MinerStateStore,
	minerConfigurationStore *store.GenericStore,
) *gin.Engine {
	allowedProxies := os.Getenv("ALLOWED_PROXIES")
	gin.SetMode(GetGinEnvironment())
	r := gin.Default()

	r.SetTrustedProxies(
		strings.Split(allowedProxies, ","),
	)
	genericroutes.CreateHealthCheckRoutes(r)
	dataroutes.CreateGenericDataRoutes(
		&r.RouterGroup,
		memoryStore,
		measurementStores,
	)
	spot.CreateSpotPriceRoutes(r, spotPriceApiClient)
	genericroutes.CreateStoreRoutes(r, measurementStores)

	v1.CreateV1Routes(
		r,
		minerConfigurationStore,
		minerStateStore,
		memoryStore,
		measurementStores,
	)
	return r
}

func PollCurrentPrice(
	spotPriceApiClient *spot.SpotHintaApiClient,
	spotPriceChan chan spot.SpotPrice,
) {
	var currentSpotPrice spot.SpotPrice
	for {
		spotPrice := spotPriceApiClient.GetCurrentPrice()
		if spotPrice != nil && spotPrice.DateTime != currentSpotPrice.DateTime {
			currentSpotPrice = *spotPrice
			spotPriceChan <- *spotPrice
		}
		time.Sleep(time.Minute)
	}
}

func WakeMiner(macAddress string) {
	log.Printf("Waking miner with MAC address %s", macAddress)
	err := exec.Command("wakeonlan", macAddress).Run()
	if err != nil {
		log.Printf("Failed to wake miner with MAC address %s: %v", macAddress, err)
	} else {
		log.Printf("Successfully sent WOL packet to miner with MAC address %s", macAddress)
	}
}

func UpdateMinerStates(
	minerStateStore *store.MinerStateStore,
	spotPriceChan chan spot.SpotPrice,
	localConfig *configuration.LocalConfig,
) {
	for {
		currentSpotPrice := <-spotPriceChan
		isMining := currentSpotPrice.PriceNoTax < localConfig.MaxSpotPriceForMining
		log.Printf("Current spot price: %f, mining: %t", currentSpotPrice.PriceNoTax, isMining)

		for _, minerId := range minerStateStore.GetIds() {
			minerState, err := minerStateStore.GetValue(minerId)
			if err != nil {
				log.Printf("Error getting miner state: %v", err)
				continue
			}

			localMinerConfig := localConfig.GetMinerConfig(minerId)
			if isMining && localMinerConfig != nil && localMinerConfig.WakeOnLan && localMinerConfig.MacAddress != "" {
				WakeMiner(localMinerConfig.MacAddress)
			}

			minerState.SpotPriceLimit = localConfig.MaxSpotPriceForMining
			minerState.TemperatureLimit = localConfig.MaxTemperatureForMining
			minerStateStore.SetValue(minerId, minerState)
		}
	}
}

func SetMinerStates(localConfig *configuration.LocalConfig, minerStateStore *store.MinerStateStore) {
	minerStateStore.Clear()
	for _, minerLocalConfig := range localConfig.Miners {
		minerState := store.MinerState{
			DeviceId:         minerLocalConfig.MinerId,
			SpotPriceLimit:   localConfig.MaxSpotPriceForMining,
			TemperatureLimit: localConfig.MaxTemperatureForMining,
		}
		minerStateStore.SetValue(minerLocalConfig.MinerId, minerState)
	}
}

func PollConfigChanges(localConfig *configuration.LocalConfig, minerStateStore *store.MinerStateStore) {
	for {
		changed := localConfig.ReloadIfNeeded()
		if changed {
			log.Println("Local config changed, updating miner states")
			SetMinerStates(localConfig, minerStateStore)
		}

		time.Sleep(time.Minute)
	}
}

func PollXmrigConfigChanges(
	minerStateStore *store.MinerStateStore,
) {
	for {
		minerIds := minerStateStore.GetIds()
		for _, minerId := range minerIds {
			stat, err := os.Stat(fmt.Sprintf("xmrig-configs/%s/config.json", minerId))
			if err != nil {
				continue
			}
			minerState, err := minerStateStore.GetValue(minerId)
			if err != nil {
				log.Printf("Error getting miner state: %v", err)
				continue
			}
			if minerState.LastConfigChanged != stat.ModTime() {
				minerState.LastConfigChanged = stat.ModTime()
				minerStateStore.SetValue(minerId, minerState)
				log.Printf("Updated config change time for miner %s", minerId)
			}
		}
		time.Sleep(time.Minute)
	}
}

func main() {
	logFile, err := os.OpenFile(GetLogFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
		return
	}
	log.SetOutput(logFile)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
		return
	}

	log.Printf("Starting %s server...", os.Getenv("ENVIRONMENT"))
	log.Println("Version: ", Version)

	var requiredEnvironmentVariables = []string{
		"API_TOKEN",
		"PORT",
	}

	if err := configvalidators.IsConfigEnvironmentVariablesValid(
		requiredEnvironmentVariables,
	); err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	var influxDbClient = store.NewInfluxDBClient()
	var memoryStore = store.MemoryStore{
		Data: make(map[string]map[string]store.Measurement),
	}
	var measurementStores = []store.MeasurementStore{}
	if influxDbClient != nil {
		measurementStores = append(measurementStores, influxDbClient)
	}

	minerStateStore := store.CreateMinerStateStore()

	spotPriceApiClient := spot.CreateSpotHintaApiClient()
	spotPriceChan := make(chan spot.SpotPrice, 1)
	localConfig := configuration.CreateLocalConfig()
	SetMinerStates(&localConfig, &minerStateStore)

	go spotPriceApiClient.PollPrices()
	go PollCurrentPrice(&spotPriceApiClient, spotPriceChan)
	go UpdateMinerStates(&minerStateStore, spotPriceChan, &localConfig)
	go PollConfigChanges(&localConfig, &minerStateStore)
	go PollXmrigConfigChanges(&minerStateStore)

	minerConfigurationStore := store.CreateGenericStore("miners_config.json")

	r := CreateRoutes(
		&memoryStore,
		measurementStores,
		&spotPriceApiClient,
		&minerStateStore,
		&minerConfigurationStore,
	)

	port := os.Getenv("PORT")

	log.Println("Server starting on port ", port)
	r.Run(":" + port)
}
