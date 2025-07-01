package miners

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"stormaaja/go-ha/data-store/store"

	"stormaaja/go-ha/data-store/assert"

	"github.com/gin-gonic/gin"
)

func TestCreateMinersRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GET /miners/:id/config - Success", func(t *testing.T) {
		router := gin.Default()
		mockConfigStore := &store.GenericStore{
			Values: map[string]any{
				"1": map[string]any{
					"config": "test-config",
				},
			},
		}
		mockStateStore := &store.MinerStateStore{}
		CreateMinersRoutes(router.Group("/v1"), mockConfigStore, mockStateStore)

		req, _ := http.NewRequest(http.MethodGet, "/v1/miners/1/config", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"config":"test-config"}`, resp.Body.String())
	})

	t.Run("GET /miners/:id/config - Miner Not Found", func(t *testing.T) {
		router := gin.Default()
		mockConfigStore := &store.GenericStore{}
		mockStateStore := &store.MinerStateStore{}
		CreateMinersRoutes(router.Group("/v1"), mockConfigStore, mockStateStore)

		req, _ := http.NewRequest(http.MethodGet, "/v1/miners/2/config", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Miner not found"}`, resp.Body.String())
	})

	t.Run("GET /miners/:id/state - Success", func(t *testing.T) {
		router := gin.Default()
		mockStateStore := &store.MinerStateStore{
			MinerStates: new(sync.Map),
		}
		mockStateStore.MinerStates.Store(
			"1",
			store.MinerState{
				DeviceId:         "1",
				IsMining:         true,
				SpotPriceLimit:   0.01,
				TemperatureLimit: 21.0,
			},
		)

		mockConfigStore := &store.GenericStore{}
		CreateMinersRoutes(router.Group("/v1"), mockConfigStore, mockStateStore)

		req, _ := http.NewRequest(http.MethodGet, "/v1/miners/1/state", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"DeviceId": "1", "IsMining": true, "LastConfigChanged":"0001-01-01T00:00:00Z", "SpotPriceLimit":0.01, "TemperatureLimit":21, "TemperatureSensorId": ""}`, resp.Body.String())
	})

	t.Run("GET /miners/:id/state - Miner Not Found", func(t *testing.T) {
		router := gin.Default()
		mockStateStore := &store.MinerStateStore{
			MinerStates: new(sync.Map),
		}
		mockConfigStore := &store.GenericStore{}
		CreateMinersRoutes(router.Group("/v1"), mockConfigStore, mockStateStore)

		req, _ := http.NewRequest(http.MethodGet, "/v1/miners/2/state", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"Miner not found"}`, resp.Body.String())
	})
}
