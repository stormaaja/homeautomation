package v1

import (
	"stormaaja/go-ha/data-store/dataroutes"
	"stormaaja/go-ha/data-store/routes/v1/shelly"
	"stormaaja/go-ha/data-store/store"

	"github.com/gin-gonic/gin"
)

func CreateV1Routes(
	g *gin.Engine,
	memoryStore *store.MemoryStore,
	measurementStores []store.MeasurementStore,
) {
	group := g.Group("/v1")
	{
		shelly.CreateShellyRoutes(group, memoryStore, measurementStores)
		dataroutes.CreateGenericDataRoutes(
			group,
			memoryStore,
			measurementStores,
		)
	}
}
