package main

import (
	"log"
	"os"
	"os/exec"
	"stormaaja/go-ha/data-store/configuration"
	"stormaaja/go-ha/data-store/configvalidators"
	"stormaaja/go-ha/data-store/dataroutes"
	"stormaaja/go-ha/data-store/genericroutes"
	v1 "stormaaja/go-ha/data-store/routes/v1"
	"stormaaja/go-ha/data-store/spot"
	"stormaaja/go-ha/data-store/state"
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

func CreateRoutes(
	memoryStore store.DataStore,
	measurementStores []store.MeasurementStore,
	spotPriceApiClient *spot.SpotHintaApiClient,
	minerStateStore *store.GenericStore,
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
		r,
		memoryStore,
		measurementStores,
	)
	spot.CreateSpotPriceRoutes(r, spotPriceApiClient)
	genericroutes.CreateStoreRoutes(r, measurementStores)

	v1.CreateV1Routes(
		r,
		minerStateStore,
		minerConfigurationStore,
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
	err := exec.Command("wakeonlan", macAddress).Run()
	if err != nil {
		log.Printf("Failed to wake miner with MAC address %s: %v", macAddress, err)
	} else {
		log.Printf("Successfully sent WOL packet to miner with MAC address %s", macAddress)
	}
}

func UpdateMinerStates(
	minerStateStore *store.GenericStore,
	spotPriceChan chan spot.SpotPrice,
	localConfig *configuration.LocalConfig,
) {
	for {
		currentSpotPrice := <-spotPriceChan

		for minerId := range minerStateStore.Values {
			minerState, err := minerStateStore.GetValue(minerId)
			if err != nil {
				log.Printf("Error getting miner state: %v", err)
				continue
			}
			minerStateMap := minerState.(map[string]any)
			isMining := currentSpotPrice.PriceNoTax < localConfig.MaxSpotPriceForMining
			localMinerConfig := localConfig.GetMinerConfig(minerId)
			if minerStateMap["isMining"] != nil && minerStateMap["isMining"].(bool) != isMining && localMinerConfig != nil && localMinerConfig.WakeOnLan && localMinerConfig.MacAddress != "" {
				WakeMiner(localMinerConfig.MacAddress)
			}

			minerStateMap["isMining"] = isMining
			minerStateStore.SetValue(minerId, minerStateMap)
		}
	}
}

func PollConfigChanges(localConfig *configuration.LocalConfig, minerStateStore *store.GenericStore) {
	for {
		changed := localConfig.ReloadIfNeeded()
		if changed {
			log.Println("Local config changed, updating miner states")
			minerIds := minerStateStore.GetIds()
			lastState := state.MinerState{IsMining: true}
			if len(minerIds) > 0 {
				lastStateFromStore, err := minerStateStore.GetValue(minerIds[0])
				if err != nil {
					log.Printf("Error getting last state: %v", err)
				} else {
					lastMinerState, ok := lastStateFromStore.(state.MinerState)
					if ok {
						lastState = lastMinerState
					} else {
						log.Printf("Error casting last state: %v", lastStateFromStore)
					}
				}
			}
			minerStateStore.Clear()
			for _, minerLocalConfig := range localConfig.Miners {
				minerState := state.MinerState{
					DeviceId: minerLocalConfig.MinerId,
					IsMining: lastState.IsMining,
				}
				minerStateStore.SetValue(minerLocalConfig.MinerId, minerState)
			}
		}

		time.Sleep(time.Minute)
	}
}

func main() {
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
		Data: make(map[string]map[string]interface{}),
	}
	var measurementStores = []store.MeasurementStore{}
	if influxDbClient != nil {
		measurementStores = append(measurementStores, influxDbClient)
	}

	minerStateStore := store.CreateGenericStore("miners_state.json")

	spotPriceApiClient := spot.CreateSpotHintaApiClient()
	spotPriceChan := make(chan spot.SpotPrice, 1)
	localConfig := configuration.CreateLocalConfig()

	go spotPriceApiClient.PollPrices()
	go PollCurrentPrice(&spotPriceApiClient, spotPriceChan)
	go UpdateMinerStates(&minerStateStore, spotPriceChan, &localConfig)
	go PollConfigChanges(&localConfig, &minerStateStore)

	minerConfigurationStore := store.CreateGenericStore("miners_config.json")

	r := CreateRoutes(
		&memoryStore,
		measurementStores,
		&spotPriceApiClient,
		&minerConfigurationStore,
		&minerStateStore,
	)

	port := os.Getenv("PORT")

	log.Println("Server starting on port ", port)
	r.Run(":" + port)
}
