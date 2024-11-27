package assert

import "testing"

func True(
	t *testing.T,
	condition bool,
) {
	if !condition {
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
	if expected != actual {
		t.Errorf(
			"Assertion failed. Expected %v, got %v",
			expected,
			actual,
		)
	}
}
