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
		group.GET("/current", func(c *gin.Context) {
			currentPrice := spotPriceApiClient.GetCurrentPrice()
			if currentPrice == nil {
				c.String(http.StatusInternalServerError, "Failed to get current price")
				return
			}
			if c.Query("format") == "priceOnly" {
				c.String(http.StatusOK, "%f", currentPrice.PriceWithTax)
				return
			}
			c.JSON(http.StatusOK, currentPrice)
		})
	}
}
