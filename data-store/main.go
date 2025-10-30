package main

import (
	"fmt"
	"log"
	"os"
	"stormaaja/go-ha/data-store/configvalidators"
	"stormaaja/go-ha/data-store/dataroutes"
	"stormaaja/go-ha/data-store/genericroutes"
	"stormaaja/go-ha/data-store/mqttclient"
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
		Data:          make(map[string]map[string]store.Measurement),
		BackupEnabled: os.Getenv("ENVIRONMENT") == "production",
	}
	if memoryStore.BackupEnabled {
		memoryStore.LoadMemoryStore()
	}
	var measurementStores = []store.MeasurementStore{}
	if influxDbClient != nil {
		measurementStores = append(measurementStores, influxDbClient)
	}

	spotPriceApiClient := spot.CreateSpotHintaApiClient()

	go spotPriceApiClient.PollPrices()

	mqttclient.Subscribe(os.Getenv("MQTT_CLIENT_ID"), os.Getenv("MQTT_BROKER"), os.Getenv("MQTT_TOPIC"), &memoryStore)

	r := CreateRoutes(
		&memoryStore,
		measurementStores,
		&spotPriceApiClient,
	)

	port := os.Getenv("PORT")

	log.Println("Server starting on port ", port)
	r.Run(":" + port)
}
