package dataroutes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"stormaaja/go-ha/data-store/store"

	"github.com/gin-gonic/gin"
)

func CreateTemperatureRoutes(
	g *gin.Engine,
	store store.DataStore,
	measurementStores []store.MeasurementStore,
) {
	g.GET("/data/temperature/:id/temperature", func(c *gin.Context) {
		sensorId := c.Param("id")
		temperature, success := store.GetFloat(sensorId)
		if !success {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("sensor not found"))
			return
		}
		c.String(http.StatusOK, "%f", temperature)
	})

	g.POST("/data/temperature/:id/temperature", func(c *gin.Context) {
		sensorId := c.Param("id")
		var temperature float64
		var body, error = io.ReadAll(c.Request.Body)
		if error != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid body"))
			return
		}
		bodyStr := string(body)
		_, error = fmt.Sscanf(bodyStr, "%f", &temperature)

		if error != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid temperature format"))
			return
		}

		store.SetFloat(sensorId, temperature)
		for _, store := range measurementStores {
			log.Printf("Storing temperature %f for sensor %s", temperature, sensorId)
			store.AppendItem("temperature", sensorId, "temperature", temperature)
		}
		c.Status(http.StatusCreated)
	})
}
