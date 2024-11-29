package middleware

import (
	"net/http"
	"net/http/httptest"
	"stormaaja/go-ha/data-store/assert"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIncludes(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"temperature", true},
		{"electricity_consumption", true},
		{"humidity", false},
		{"", false},
	}

	for _, test := range tests {
		result := Includes(test.value)
		assert.Equal(t, test.expected, result)
	}
}

func TestMeasurementTypeValidator(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		measurementType string
		expectedStatus  int
	}{
		{"temperature", http.StatusOK},
		{"electricity_consumption", http.StatusOK},
		{"humidity", http.StatusBadRequest},
		{"", http.StatusBadRequest},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{gin.Param{Key: "measurement", Value: test.measurementType}}
		MeasurementTypeValidator()(c)

		assert.Equal(t, test.expectedStatus, w.Code)
	}
}
