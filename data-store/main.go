package main

import (
	"log"
	"os"
	"stormaaja/go-ha/data-store/dataroutes"
	"stormaaja/go-ha/data-store/genericroutes"
	"stormaaja/go-ha/data-store/middleware"
	"stormaaja/go-ha/data-store/store"
	"stormaaja/go-ha/security/configvalidators"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CreateRoutes(
	memoryStore store.DataStore,
	measurementStores []store.MeasurementStore,
) *gin.Engine {
	allowedProxies := os.Getenv("ALLOWED_PROXIES")
	r := gin.Default()
	switch os.Getenv("ENVIRONMENT") {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

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
	r := CreateRoutes(
		&memoryStore,
		[]store.MeasurementStore{&influxDbClient},
	)
	port := os.Getenv("PORT")

	log.Println("Server starting on port ", port)
	r.Run(":" + port)
}
