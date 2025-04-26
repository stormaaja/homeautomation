package assert

import (
	"encoding/json"
	"maps"
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

func Equal(
	t *testing.T,
	expected interface{},
	actual interface{},
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
	parsedTarget := map[string]interface{}{}
	parsedActual := map[string]interface{}{}
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
