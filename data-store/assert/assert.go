package assert

import "testing"

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
