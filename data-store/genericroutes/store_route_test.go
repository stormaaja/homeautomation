package genericroutes

import (
	"net/http"
	"net/http/httptest"
	"stormaaja/go-ha/data-store/store"
	"testing"
)

type MockMeasurementStore struct {
	Data []float64
}

func (m *MockMeasurementStore) AppendItem(measurement string, location string, field string, value float64) {
	m.Data = append(m.Data, value)
}

func (m *MockMeasurementStore) Flush() {
	m.Data = []float64{}
}

func TestHandleGet(t *testing.T) {
	storeRoute := StoreRoute{}
	req, err := http.NewRequest("GET", "/measurements", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(storeRoute.HandleGet)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHandlePostInvalidPath(t *testing.T) {
	storeRoute := StoreRoute{}
	req, err := http.NewRequest("POST", "/measurements/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(storeRoute.HandlePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandlePostFlush(t *testing.T) {
	mockStore := MockMeasurementStore{}
	storeRoute := StoreRoute{MeasurementStores: []store.MeasurementStore{&mockStore}}
	mockStore.AppendItem("test-measurement", "test-sensor", "test-field", 1.0)

	if len(mockStore.Data) != 1 {
		t.Errorf("Data not appended")
	}

	req, err := http.NewRequest("POST", "/measurements/flush", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(storeRoute.HandlePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v (%v)", status, http.StatusOK, rr.Body.String())
	}

	if len(mockStore.Data) != 0 {
		t.Errorf("Data not flushed")
	}

}
