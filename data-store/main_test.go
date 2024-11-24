package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"stormaaja/go-ha/data-store/store"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("API_TOKEN", "valid-token")
	memoryStore := store.MemoryStore{Data: make(map[string]interface{})}

	router := CreateRoutes(
		&memoryStore,
		[]store.MeasurementStore{},
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Expected status code 401, got %v", w.Code)
	}
}

func TestValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("API_TOKEN", "valid-token")
	memoryStore := store.MemoryStore{Data: make(map[string]interface{})}

	router := CreateRoutes(
		&memoryStore,
		[]store.MeasurementStore{},
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}

}
