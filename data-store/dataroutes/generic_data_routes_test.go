package dataroutes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"stormaaja/go-ha/data-store/assert"
	"stormaaja/go-ha/data-store/store"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateGenericDataRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDataStore := store.MemoryStore{
		Data: make(map[string]map[string]interface{}),
	}
	mockMeasurementStore := store.MockMeasurementStore{
		Items: make(map[string]float64),
	}
	measurementStores := []store.MeasurementStore{&mockMeasurementStore}
	CreateGenericDataRoutes(router, mockDataStore, measurementStores)

	t.Run("GET /data/:measurement/:id/:field - success", func(t *testing.T) {
		mockDataStore.SetMeasurement("testtype", "device1", store.Measurement{
			DeviceId:        "device1",
			MeasurementType: "testtype",
			Field:           "temperature",
			Value:           25.5,
		})

		req, _ := http.NewRequest(http.MethodGet, "/data/testtype/device1/temperature", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "25.500000", resp.Body.String())
	})

	t.Run("GET /data/:measurement/:id/:field - device not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/data/temperature/device2/temperature", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("POST /data/:measurement/:id/:field - success", func(t *testing.T) {
		body := bytes.NewBufferString("30.5")
		req, _ := http.NewRequest(http.MethodPost, "/data/testtype/device2/testfield", body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, 30.5, mockMeasurementStore.Items["device2"])
		assert.Equal(t, 30.5, mockDataStore.Data["device2"]["testtype"].(store.Measurement).Value)
	})

	t.Run("POST /data/:measurement/:id/:field - invalid body", func(t *testing.T) {
		body := bytes.NewBufferString("invalid")
		req, _ := http.NewRequest(http.MethodPost, "/data/temperature/device1/temperature", body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}