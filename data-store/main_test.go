package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"stormaaja/go-ha/data-store/spot"
	"stormaaja/go-ha/data-store/store"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetGinEnvironment(t *testing.T) {
	tests := []struct {
		env      string
		expected string
	}{
		{"production", gin.ReleaseMode},
		{"test", gin.TestMode},
		{"development", gin.DebugMode},
	}

	for _, test := range tests {
		os.Setenv("ENVIRONMENT", test.env)
		result := GetGinEnvironment()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("API_TOKEN", "valid-token")
	memoryStore := store.MemoryStore{Data: make(map[string]map[string]any)}

	router := CreateRoutes(
		&memoryStore,
		[]store.MeasurementStore{},
		&spot.SpotHintaApiClient{},
		&store.MinerStateStore{},
		&store.GenericStore{},
	)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString("30.5")
	req, _ := http.NewRequest(http.MethodPost, "/data/temperature/device2/temperature", body)
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Expected status code 401, got %v", w.Code)
	}
}

func TestValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("API_TOKEN", "valid-token")
	memoryStore := store.MemoryStore{Data: make(map[string]map[string]any)}
	memoryStore.SetMeasurement("temperature", "device-id", store.Measurement{
		DeviceId:        "device-id",
		MeasurementType: "temperature",
		Field:           "temperature",
		Value:           25.5,
	})

	router := CreateRoutes(
		&memoryStore,
		[]store.MeasurementStore{},
		&spot.SpotHintaApiClient{},
		&store.MinerStateStore{},
		&store.GenericStore{},
	)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString("30.5")
	req, _ := http.NewRequest(http.MethodPost, "/data/temperature/device2/temperature", body)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected status code 201, got %v", w.Code)
	}

}
