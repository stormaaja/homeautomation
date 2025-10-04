package dataroutes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"stormaaja/go-ha/data-store/assert"
	"stormaaja/go-ha/data-store/store"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateGenericDataRoutes(t *testing.T) {
	os.Setenv("API_TOKEN", "valid-token")
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	memoryStore := store.MemoryStore{
		Data: make(map[string]map[string]any),
	}
	mockMeasurementStore := store.MockMeasurementStore{
		Items: make(map[string]float64),
	}
	measurementStores := []store.MeasurementStore{&mockMeasurementStore}
	CreateGenericDataRoutes(&router.RouterGroup, &memoryStore, measurementStores)

	t.Run("GET /data/:measurement/:id/:field - success", func(t *testing.T) {
		memoryStore.SetMeasurement("electricity_consumption", "device1", store.Measurement{
			DeviceId:        "device1",
			MeasurementType: "electricity_consumption",
			Field:           "energy",
			Value:           25.5,
		})

		req, _ := http.NewRequest(http.MethodGet, "/data/electricity_consumption/device1/energy", nil)
		req.Header.Add("Authorization", "Bearer valid-token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "25.500000", resp.Body.String())
	})

	t.Run("GET /data/:measurement/:id/:field?format=full - success", func(t *testing.T) {
		memoryStore.SetMeasurement("electricity_consumption", "device1", store.Measurement{
			DeviceId:        "device1",
			MeasurementType: "electricity_consumption",
			Field:           "energy",
			Value:           25.5,
		})

		req, _ := http.NewRequest(http.MethodGet, "/data/electricity_consumption/device1/energy?format=full", nil)
		req.Header.Add("Authorization", "Bearer valid-token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, "{\"DeviceId\":\"device1\",\"MeasurementType\":\"electricity_consumption\",\"Field\":\"energy\",\"Value\":25.5,\"UpdatedAt\":\"0001-01-01T00:00:00Z\"}", resp.Body.String())
	})

	t.Run("GET /data/:measurement/:id/:field - device not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/data/temperature/device2/temperature", nil)
		req.Header.Add("Authorization", "Bearer valid-token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("POST /data/:measurement/:id/:field - success", func(t *testing.T) {
		body := bytes.NewBufferString("30.5")
		req, _ := http.NewRequest(http.MethodPost, "/data/temperature/device2/temperature", body)
		req.Header.Add("Authorization", "Bearer valid-token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, 30.5, mockMeasurementStore.Items["device2"])
		assert.Equal(t, 30.5, memoryStore.Data["device2"]["temperature"].(store.Measurement).Value)
	})

	t.Run("POST /data/:measurement/:id/:field - invalid body", func(t *testing.T) {
		body := bytes.NewBufferString("invalid")
		req, _ := http.NewRequest(http.MethodPost, "/data/temperature/device1/temperature", body)
		req.Header.Add("Authorization", "Bearer valid-token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
