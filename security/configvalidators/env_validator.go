package configvalidators

import (
	"errors"
	"os"
)

// IsConfigEnvironmentVariablesValid validates environment variables
//
// It returns error, if the validation is failed. Otherwise, it returns nil.
func IsConfigEnvironmentVariablesValid() error {
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		return errors.New("API_TOKEN is not set")
	}
	return nil
}
