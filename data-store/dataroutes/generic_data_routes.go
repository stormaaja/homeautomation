package dataroutes

import (
	"fmt"
	"log"
	"net/http"
	"stormaaja/go-ha/data-store/store"
	"stormaaja/go-ha/data-store/tools"

	"github.com/gin-gonic/gin"
)

func CreateGenericDataRoutes(
	g *gin.Engine,
	datastore store.DataStore,
	measurementStores []store.MeasurementStore,
) {
	g.GET("/data/:measurement/:id/:field", func(c *gin.Context) {
		measurementType := c.Param("measurement")
		deviceId := c.Param("id")
		field := c.Param("field")
		measurement, success := datastore.GetMeasurement(measurementType, deviceId)
		if !success {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("device not found"))
			return
		}
		var valueString string = ""
		switch field {
		case "temperature", "energy":
			valueString = fmt.Sprintf("%f", measurement.Value.(float64))
		default:
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("field not found"))
			return
		}
		c.String(http.StatusOK, valueString)
	})

	g.POST("/data/:measurement/:id/:field", func(c *gin.Context) {
		measurementType := c.Param("measurement")
		deviceId := c.Param("id")
		field := c.Param("field")
		value, error := tools.ReadBodyFloat(&c.Request.Body)

		if error != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid body"))
			return
		}

		measurement := store.Measurement{
			DeviceId:        deviceId,
			MeasurementType: measurementType,
			Field:           field,
			Value:           value,
		}

		datastore.SetMeasurement(
			measurementType,
			deviceId,
			measurement,
		)

		for _, measurementStore := range measurementStores {
			log.Printf("Storing value %f for device %s", measurement.Value, measurement.DeviceId)
			measurementStore.AppendItem(
				measurementType,
				deviceId,
				field,
				value,
			)
		}
		c.Status(http.StatusCreated)
	})
}