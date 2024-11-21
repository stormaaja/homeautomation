package requestvalidators

import (
	"net/http"
	"os"
	"testing"
)

func TestValidateTokenValidToken(t *testing.T) {
	// Set up the environment variable for testing
	validToken := "valid-token"
	os.Setenv("API_TOKEN", validToken)
	defer os.Unsetenv("API_TOKEN")
	header := http.Header{
		"Authorization": []string{"Bearer " + validToken},
	}

	if !IsApiTokenValid(header) {
		t.Errorf("ValidateToken() returned false")
	}

}
