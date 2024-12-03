package configvalidators

import (
	"os"
	"testing"
)

func TestIsConfigEnvironmentVariablesValid(t *testing.T) {
	t.Run("API_TOKEN is set", func(t *testing.T) {
		os.Setenv("API_TOKEN", "some-token")
		defer os.Unsetenv("API_TOKEN")

		err := IsConfigEnvironmentVariablesValid([]string{"API_TOKEN"})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("API_TOKEN is not set", func(t *testing.T) {
		os.Unsetenv("API_TOKEN")

		err := IsConfigEnvironmentVariablesValid([]string{"API_TOKEN"})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
