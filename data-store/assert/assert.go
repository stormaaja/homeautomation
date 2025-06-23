package assert

import (
	"encoding/json"
	"maps"
	"strings"
	"testing"
)

func True(
	t *testing.T,
	condition bool,
) {
	if !condition {
		t.Helper()
		t.Errorf(
			"Assertion failed. Expected true, got false. Message",
		)
	}
}

func False(
	t *testing.T,
	condition bool,
) {
	if condition {
		t.Helper()
		t.Errorf(
			"Assertion failed. Expected false, got true. Message",
		)
	}
}

func Equal(
	t *testing.T,
	expected any,
	actual any,
) {
	t.Helper()
	if expected != actual {
		t.Errorf(
			"Assertion failed. Expected %v, got %v",
			expected,
			actual,
		)
	}
}

func JSONEq(
	t *testing.T,
	expected string,
	actual string,
) {
	t.Helper()
	parsedTarget := map[string]any{}
	parsedActual := map[string]any{}
	err := json.Unmarshal([]byte(expected), &parsedTarget)
	if err != nil {
		t.Errorf(
			"Failed to parse expected JSON: %v",
			err,
		)
		return
	}
	err = json.Unmarshal([]byte(actual), &parsedActual)
	if err != nil {
		t.Errorf(
			"Failed to parse actual JSON: %v",
			err,
		)
		return
	}
	if !maps.Equal(parsedTarget, parsedActual) {
		t.Errorf(
			"Assertion failed. Expected JSON %s, got %s",
			expected,
			actual,
		)
	}
}

func Contains(
	t *testing.T,
	container string,
	substring string,
) {
	t.Helper()
	if !strings.Contains(container, substring) {
		t.Errorf(
			"Assertion failed. Expected '%s' to contain '%s'",
			container,
			substring,
		)
	}
}
