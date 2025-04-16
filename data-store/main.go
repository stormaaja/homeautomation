package main

import (
	"log"
	"os"
	"stormaaja/go-ha/data-store/configvalidators"
	"stormaaja/go-ha/data-store/dataroutes"
	"stormaaja/go-ha/data-store/genericroutes"
	"stormaaja/go-ha/data-store/middleware"
	"stormaaja/go-ha/data-store/store"

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
) *gin.Engine {
	allowedProxies := os.Getenv("ALLOWED_PROXIES")
	gin.SetMode(GetGinEnvironment())
	r := gin.Default()

	r.SetTrustedProxies(
		strings.Split(allowedProxies, ","),
	)
	r.Use(middleware.TokenCheck())
	genericroutes.CreateHealthCheckRoutes(r)
	dataroutes.CreateGenericDataRoutes(
		r,
		memoryStore,
		measurementStores,
	)
	genericroutes.CreateStoreRoutes(r, measurementStores)
	return r
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
	r := CreateRoutes(
		&memoryStore,
		measurementStores,
	)

	port := os.Getenv("PORT")

	log.Println("Server starting on port ", port)
	r.Run(":" + port)
}
