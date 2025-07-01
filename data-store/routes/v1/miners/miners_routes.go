package miners

import (
	"fmt"
	"log"
	"net/http"
	"stormaaja/go-ha/data-store/middleware"
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
		minerGroup := group.Group("/:id")
		{
			minerGroup.GET("/config", middleware.Deprecated(), func(c *gin.Context) {

				id := c.Param("id")
				minerConfig, err := configurationStore.GetValue(id)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Miner not found"})
					return
				}
				c.JSON(http.StatusOK, minerConfig)
			})
			minerGroup.GET("/state", func(c *gin.Context) {
				id := c.Param("id")
				minerState, err := stateStore.GetValue(id)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Miner not found"})
					return
				}
				c.JSON(http.StatusOK, minerState)
			})
			minerGroup.GET("/xmrig/config.json", func(c *gin.Context) {
				id := c.Param("id")
				c.File(fmt.Sprintf("xmrig-configs/%s/config.json", id))
			})
		}
		changeGroup := group.Group("/refresh")
		{
			changeGroup.Use(middleware.TokenCheck())
			changeGroup.POST("", func(c *gin.Context) {
				log.Println("Refreshing miner state and configuration")
				stateStore.Load()
				configurationStore.Load()
				c.JSON(http.StatusOK, gin.H{"status": "success"})
			})
		}
	}
}
