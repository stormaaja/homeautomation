package dataroutes

import (
	"fmt"
	"log"
	"net/http"
	"stormaaja/go-ha/data-store/store"
	"stormaaja/go-ha/data-store/tools"

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
		temperature, error := tools.ReadBodyFloat(&c.Request.Body)

		if error != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid body"))
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
