package dataroutes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"stormaaja/go-ha/data-store/store"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockMeasurementStore struct {
	Items map[string]float64
}

func (m *MockMeasurementStore) AppendItem(
	measurement string,
	location string,
	field string,
	value float64,
) {
	m.Items[location] = value
}

func (m *MockMeasurementStore) Flush() {
	m.Items = make(map[string]float64)
}

func TestGetTemperature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockStore := store.MemoryStore{Data: make(map[string]interface{})}
	mockStore.SetFloat("sensor1", 25.5)

	CreateTemperatureRoutes(router, &mockStore, nil)

	req, _ := http.NewRequest(http.MethodGet, "/data/temperature/sensor1/temperature", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if http.StatusOK != resp.Code {
		t.Errorf("Expected status code 200, got %v", resp.Code)
	}

	if resp.Body.String() != "25.500000" {
		t.Errorf("Expected 25.500000, got %v", resp.Body.String())
	}
}

func TestGetTemperatureNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockStore := store.MemoryStore{Data: make(map[string]interface{})}

	CreateTemperatureRoutes(router, &mockStore, nil)

	req, _ := http.NewRequest(http.MethodGet, "/data/temperature/sensor1/temperature", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if http.StatusBadRequest != resp.Code {
		t.Errorf("Expected status code 400, got %v", resp.Code)
	}
}

func TestPostTemperature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockStore := store.MemoryStore{Data: make(map[string]interface{})}
	mockMeasurementStore := MockMeasurementStore{Items: make(map[string]float64)}
	measurementStores := []store.MeasurementStore{&mockMeasurementStore}

	CreateTemperatureRoutes(router, &mockStore, measurementStores)

	body := bytes.NewBufferString("30.5")
	req, _ := http.NewRequest(http.MethodPost, "/data/temperature/sensor1/temperature", body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if http.StatusCreated != resp.Code {
		t.Errorf("Expected status code 201, got %v", resp.Code)
	}

	if mockStore.Data["sensor1"] != 30.5 {
		t.Errorf("Expected 30.5, got %v", mockStore.Data["sensor1"])
	}

	if mockMeasurementStore.Items["sensor1"] != 30.5 {
		t.Errorf("Expected 30.5, got %v", mockMeasurementStore.Items["sensor1"])
	}
}

func TestPostTemperatureInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockStore := store.MemoryStore{Data: make(map[string]interface{})}

	CreateTemperatureRoutes(router, &mockStore, nil)

	body := bytes.NewBufferString("invalid")
	req, _ := http.NewRequest(http.MethodPost, "/data/temperature/sensor1/temperature", body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if http.StatusBadRequest != resp.Code {
		t.Errorf("Expected status code 400, got %v", resp.Code)
	}
}
