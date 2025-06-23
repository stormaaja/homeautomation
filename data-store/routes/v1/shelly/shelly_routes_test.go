package shelly

import (
	"net/http"
	"net/http/httptest"
	"stormaaja/go-ha/data-store/assert"
	"stormaaja/go-ha/data-store/store"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestReportValues_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	memoryStore := store.MemoryStore{
		Data: make(map[string]map[string]any),
	}
	measurementStore := store.MockMeasurementStore{
		Items: make(map[string]float64),
	}

	CreateShellyRoutes(router.Group("/v1"), &memoryStore, []store.MeasurementStore{&measurementStore})

	req, _ := http.NewRequest("GET", "/v1/shelly/ht/test-device/report-values?temp=23.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	measurement := memoryStore.Data["test-device"]["temperature"].(store.Measurement)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "test-device", measurement.DeviceId)
	assert.Equal(t, 23.5, measurement.Value)
}

func TestReportValues_MissingTemperature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	memoryStore := store.MemoryStore{
		Data: make(map[string]map[string]any),
	}
	measurementStore := store.MockMeasurementStore{
		Items: make(map[string]float64),
	}

	CreateShellyRoutes(router.Group("/v1"), &memoryStore, []store.MeasurementStore{&measurementStore})

	req, _ := http.NewRequest("GET", "/v1/shelly/ht/test-device/report-values", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing temperature")
}

func TestReportValues_InvalidTemperature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	memoryStore := store.MemoryStore{
		Data: make(map[string]map[string]any),
	}
	measurementStore := store.MockMeasurementStore{
		Items: make(map[string]float64),
	}

	CreateShellyRoutes(router.Group("/v1"), &memoryStore, []store.MeasurementStore{&measurementStore})

	req, _ := http.NewRequest("GET", "/v1/shelly/ht/test-device/report-values?temp=notanumber", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid value for temperature: notanumber")
}
