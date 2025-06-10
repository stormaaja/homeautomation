package shelly

import (
	"log"
	"net/http"
	"stormaaja/go-ha/data-store/store"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateShellyRoutes(
	g *gin.RouterGroup,
	datastore store.DataStore,
	measurementStores []store.MeasurementStore,
) {
	// /v1/shelly/ht/shellyhtdownstairs1?hum=36&temp=25.12&id=shellyht-3C63BB
	group := g.Group("/shelly")
	{
		htGroup := group.Group("/ht/:id")
		{
			// Shelly H&T uses GET for reporting values. That's not a good practice, but we have to deal with it.
			htGroup.GET("/report-values", func(c *gin.Context) {
				deviceId := c.Param("id")
				temperature := c.Query("temp")

				if temperature == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Missing temperature"})
					return
				}

				value, err := strconv.ParseFloat(temperature, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid temperature value"})
					return
				}

				temperatureMeasurement := store.Measurement{
					DeviceId:        deviceId,
					MeasurementType: "temperature",
					Field:           "temperature",
					Value:           value,
				}

				datastore.SetMeasurement(
					temperatureMeasurement.MeasurementType,
					deviceId,
					temperatureMeasurement,
				)

				for _, measurementStore := range measurementStores {
					log.Printf("Storing value %f for device %s", temperatureMeasurement.Value, temperatureMeasurement.DeviceId)
					measurementStore.AppendItem(
						temperatureMeasurement.MeasurementType,
						deviceId,
						temperatureMeasurement.Field,
						value,
					)
				}
				c.Status(http.StatusCreated)
			})
		}
	}
}
