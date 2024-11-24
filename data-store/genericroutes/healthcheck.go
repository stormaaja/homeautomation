package genericroutes

import (
	"github.com/gin-gonic/gin"
)

func CreateHealthCheckRoutes(g *gin.Engine) {
	g.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, "")
	})
}
