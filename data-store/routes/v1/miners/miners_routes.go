package miners

import (
	"net/http"
	"stormaaja/go-ha/data-store/store"

	"github.com/gin-gonic/gin"
)

func CreateMinersRoutes(
	g *gin.RouterGroup,
	configurationStore *store.GenericStore,
	stateStore *store.MinerStateStore,
) {
	group := g.Group("/miners")
	{
		group.GET("/:id/config", func(c *gin.Context) {
			id := c.Param("id")
			minerConfig, err := configurationStore.GetValue(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Miner not found"})
				return
			}
			c.JSON(http.StatusOK, minerConfig)
		})
		group.GET("/:id/state", func(c *gin.Context) {
			id := c.Param("id")
			minerState, err := stateStore.GetValue(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Miner not found"})
				return
			}
			c.JSON(http.StatusOK, minerState)
		})
	}
}
