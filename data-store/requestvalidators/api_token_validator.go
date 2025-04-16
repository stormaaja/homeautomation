package requestvalidators

import (
	"net/http"
	"os"
)

// IsApiTokenValid validates the token in the header
//
// It returns true if the validation is success, false otherwise
func IsApiTokenValid(header http.Header) bool {
	apiToken := os.Getenv("API_TOKEN")
	return header.Get("Authorization") == "Bearer "+apiToken
}
