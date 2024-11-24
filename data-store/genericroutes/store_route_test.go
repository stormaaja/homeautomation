package genericroutes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"stormaaja/go-ha/data-store/store"

	"github.com/gin-gonic/gin"
)

type MockMeasurementStore struct {
	Flushed bool
}

func (m *MockMeasurementStore) Flush() {
	m.Flushed = true
}

func (m *MockMeasurementStore) AppendItem(
	measurement string,
	location string,
	field string,
	value float64,
) {
}

func TestCreateStoreRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockStore1 := &MockMeasurementStore{}
	mockStore2 := &MockMeasurementStore{}
	measurementStores := []store.MeasurementStore{mockStore1, mockStore2}

	CreateStoreRoutes(router, measurementStores)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/measurements/flush", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}

	if !mockStore1.Flushed {
		t.Errorf("Expected store 1 to be flushed")
	}

	if !mockStore2.Flushed {
		t.Errorf("Expected store 2 to be flushed")
	}
}
