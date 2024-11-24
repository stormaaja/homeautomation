package genericroutes

import (
	"stormaaja/go-ha/data-store/store"

	"github.com/gin-gonic/gin"
)

func CreateStoreRoutes(g *gin.Engine, measurementStores []store.MeasurementStore) {
	g.POST("/measurements/flush", func(ctx *gin.Context) {
		for _, store := range measurementStores {
			store.Flush()
		}
	})
}
