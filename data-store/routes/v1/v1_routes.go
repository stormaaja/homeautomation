package v1

import (
	"stormaaja/go-ha/data-store/routes/v1/miners"
	"stormaaja/go-ha/data-store/store"

	"github.com/gin-gonic/gin"
)

func CreateV1Routes(
	g *gin.Engine,
	configurationStore *store.GenericStore,
	stateStore *store.GenericStore,
) {
	group := g.Group("/v1")
	{
		miners.CreateMinersRoutes(group, configurationStore, stateStore)
	}
}
