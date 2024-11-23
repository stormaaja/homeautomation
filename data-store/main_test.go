package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("API_TOKEN", "valid-token")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestHandlerWithValidToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("API_TOKEN", "valid-token")
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandlerWithoutPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("API_TOKEN", "valid-token")
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
func TestGetRootPath(t *testing.T) {
	tests := []struct {
		path         string
		expectedRoot string
		expectedSub  string
	}{
		{"/data/temperature", "data", "temperature"},
		{"/data/", "data", ""},
		{"/healthcheck", "healthcheck", ""},
		{"/", "", ""},
		{"", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			root, sub := GetRootPath(tt.path)
			if root != tt.expectedRoot || sub != tt.expectedSub {
				t.Errorf("GetRootPath(%s) = (%s, %s); want (%s, %s)", tt.path, root, sub, tt.expectedRoot, tt.expectedSub)
			}
		})
	}
}
