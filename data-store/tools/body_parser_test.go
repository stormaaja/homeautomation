package tools

import (
	"io"
	"strings"
	"testing"
)

func TestReadBodyString(t *testing.T) {
	body := "test body"
	readCloser := io.NopCloser(strings.NewReader(body))

	result, err := ReadBodyString(&readCloser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != body {
		t.Errorf("Expected %v, got %v", body, result)
	}
}

func TestReadBodyFloat(t *testing.T) {
	body := "123.45"
	readCloser := io.NopCloser(strings.NewReader(body))

	result, err := ReadBodyFloat(&readCloser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := 123.45
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReadBodyFloat_Invalid(t *testing.T) {
	body := "invalid float"
	readCloser := io.NopCloser(strings.NewReader(body))

	_, err := ReadBodyFloat(&readCloser)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
