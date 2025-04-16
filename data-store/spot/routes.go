package spot

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSpotPriceRoutes(
	g *gin.Engine,
	spotPriceApiClient *SpotHintaApiClient,
) {
	group := g.Group("/v1/spot")
	{
		group.GET("/prices", func(c *gin.Context) {
			prices := spotPriceApiClient.GetPrices()
			if prices == nil {
				c.String(http.StatusInternalServerError, "Failed to get prices")
				return
			}
			c.JSON(http.StatusOK, prices)
		})
	}
}
