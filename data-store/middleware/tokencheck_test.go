package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTokenCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Token", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		os.Setenv("API_TOKEN", "valid-token")
		r.Use(TokenCheck())
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %v", w.Code)
		}

		if w.Body.String() != "OK" {
			t.Errorf("Expected 'OK', got %v", w.Body.String())
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		os.Setenv("API_TOKEN", "valid-token")
		r.Use(TokenCheck())
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got %v", w.Code)
		}

		if w.Body.String() != "unauthorized" {
			t.Errorf("Expected 'unauthorized', got %v", w.Body.String())
		}
	})

	t.Run("Without Token", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		os.Setenv("API_TOKEN", "valid-token")
		r.Use(TokenCheck())
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got %v", w.Code)
		}

		if w.Body.String() != "unauthorized" {
			t.Errorf("Expected 'unauthorized', got %v", w.Body.String())
		}
	})
}
