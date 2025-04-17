package spot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"stormaaja/go-ha/data-store/assert"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type MockSpotHintaApiClient struct{}

func TestCreateSpotPriceRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	currentTime := time.Now()
	apiClient := &SpotHintaApiClient{
		State: SpotHintaApiState{
			SpotPrices: []SpotPrice{
				{Rank: 1, DateTime: currentTime.Truncate(time.Hour), PriceNoTax: 0.1, PriceWithTax: 0.12, PercentageRank: 0.5},
				{Rank: 2, DateTime: currentTime.Truncate(time.Hour).Add(time.Hour), PriceNoTax: 0.2, PriceWithTax: 0.24, PercentageRank: 1.0},
			},
			LastCheck: time.Now(),
		},
	}

	CreateSpotPriceRoutes(router, apiClient)

	t.Run("GET /v1/spot/prices", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/spot/prices", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		parsedResponse := []SpotPrice{}
		json.Unmarshal(resp.Body.Bytes(), &parsedResponse)
		assert.Equal(t, apiClient.State.SpotPrices[0].Rank, parsedResponse[0].Rank)
		assert.Equal(t, apiClient.State.SpotPrices[1].Rank, parsedResponse[1].Rank)
	})

	t.Run("GET /v1/spot/current", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/spot/current", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		parsedResponse := SpotPrice{}
		json.Unmarshal(resp.Body.Bytes(), &parsedResponse)
		assert.Equal(t, apiClient.State.SpotPrices[0].Rank, parsedResponse.Rank)
	})

	t.Run("GET /v1/spot/current with format=priceOnly", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/spot/current?format=priceOnly", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "0.120000", resp.Body.String())
	})
}
