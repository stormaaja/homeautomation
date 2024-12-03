package configvalidators

import (
	"errors"
	"os"
)

// IsConfigEnvironmentVariablesValid validates environment variables
//
// It returns error, if the validation is failed. Otherwise, it returns nil.
func IsConfigEnvironmentVariablesValid(requiredEnvironmentVariables []string) error {
	for _, environmentVariable := range requiredEnvironmentVariables {
		if os.Getenv(environmentVariable) == "" {
			return errors.New(environmentVariable + " is not set")
		}
	}
	return nil
}
