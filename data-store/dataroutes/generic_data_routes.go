package dataroutes

import (
	"fmt"
	"log"
	"net/http"
	"stormaaja/go-ha/data-store/middleware"
	"stormaaja/go-ha/data-store/store"
	"stormaaja/go-ha/data-store/tools"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateGenericDataRoutes(
	g *gin.RouterGroup,
	memoryStore *store.MemoryStore,
	measurementStores []store.MeasurementStore,
) {
	dataGroup := g.Group("/data")
	{
		group := dataGroup.Group("/:measurement/:id/:field")
		{
			group.Use(middleware.MeasurementTypeValidator())

			group.GET("", func(c *gin.Context) {
				measurementType := c.Param("measurement")
				deviceId := c.Param("id")
				field := c.Param("field")
				measurement, success := memoryStore.GetMeasurement(measurementType, deviceId)
				if !success {
					c.AbortWithError(http.StatusBadRequest, fmt.Errorf("device or measurement type not found"))
					return
				}
				if c.Query("format") == "full" {
					c.JSON(http.StatusOK, measurement)
				} else {
					var valueString string = ""
					switch field {
					case "temperature", "energy":
						valueString = fmt.Sprintf("%f", measurement.Value.(float64))
					default:
						c.AbortWithError(http.StatusBadRequest, fmt.Errorf("field not found"))
						return
					}
					c.String(http.StatusOK, valueString)
				}
			})

			group.POST("", middleware.TokenCheck(), func(c *gin.Context) {
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
					UpdatedAt:       time.Now(),
				}

				memoryStore.SetMeasurement(
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
	}

}
