package shelly

import (
	"fmt"
	"log"
	"net/http"
	"stormaaja/go-ha/data-store/store"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StoreTemperature(
	datastore store.DataStore,
	measurementStores []store.MeasurementStore,
	deviceId string,
	valueStr string,
	valueType string,
) error {
	if valueStr == "" {
		return fmt.Errorf("missing %s", valueType)
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("invalid value for %s: %s", valueType, valueStr)
	}

	measurement := store.Measurement{
		DeviceId:        deviceId,
		MeasurementType: valueType,
		Field:           valueType,
		Value:           value,
	}

	datastore.SetMeasurement(
		measurement.MeasurementType,
		deviceId,
		measurement,
	)

	for _, measurementStore := range measurementStores {
		log.Printf("Storing value %f for device %s", measurement.Value, measurement.DeviceId)
		measurementStore.AppendItem(
			measurement.MeasurementType,
			deviceId,
			measurement.Field,
			value,
		)
	}
	return nil
}

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

				err := StoreTemperature(datastore, measurementStores, deviceId, temperature, "temperature")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}

				c.Status(http.StatusCreated)
			})
		}
	}
}
