package genericroutes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateHealthCheckRoutes(t *testing.T) {
	router := gin.Default()
	CreateHealthCheckRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}
}
