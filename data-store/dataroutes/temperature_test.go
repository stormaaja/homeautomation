package dataroutes

import (
	"net/http"
	"net/http/httptest"
	"stormaaja/go-ha/data-store/store"
	"strings"
	"testing"
)

func TestParseId(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/data/temperature/123/temperature", "123"},
		{"/data/temperature/", ""},
		{"/data/temperature/123/extra", "123"},
		{"/", ""},
	}

	for _, test := range tests {
		result := ParseId(test.path)
		if result != test.expected {
			t.Errorf("ParseId(%s) = %s; expected %s", test.path, result, test.expected)
		}
	}
}

func TestIsValidValueType(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/data/temperature/123/temperature", true},
		{"/data/temperature/123/humidity", false},
		{"/data/temperature/123/", false},
		{"/data/temperature/123/temperature/extra", false},
		{"/data/temperature/123/temperature", true},
	}

	for _, test := range tests {
		result := IsValidValueType(test.path)
		if result != test.expected {
			t.Errorf("IsValidValueType(%s) = %v; expected %v", test.path, result, test.expected)
		}
	}
}

func TestHandleGet(t *testing.T) {
	dataStore := store.MemoryStore{Data: make(map[string]interface{})}
	dataStore.SetFloat("123", 25.5)
	route := TemperatureRoute{Store: &dataStore}

	req, err := http.NewRequest("GET", "/data/temperature/123/temperature", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(route.HandleGet)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "25.500000"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleGetNonExisting(t *testing.T) {
	dataStore := store.MemoryStore{Data: make(map[string]interface{})}
	route := TemperatureRoute{Store: &dataStore}

	req, err := http.NewRequest("GET", "/data/temperature/123/temperature", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(route.HandleGet)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandlePost(t *testing.T) {
	dataStore := store.MemoryStore{Data: make(map[string]interface{})}
	route := TemperatureRoute{Store: &dataStore}

	body := strings.NewReader("30.5")
	req, err := http.NewRequest("POST", "/date/temperature/123/temperature", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(route.HandlePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if value, _ := dataStore.GetFloat("123"); value != 30.5 {
		t.Errorf("handler did not set the correct value: got %v want %v", value, 30.5)
	}
}
