package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Helper to execute a request through the middleware
func executeRequest(t *testing.T, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	// Set a User-Agent header so /headers test passes
	req.Header.Set("User-Agent", "httptest")

	rr := httptest.NewRecorder()

	// Match route manually since we are testing handlers directly
	switch path {
	case "/":
		loggingMiddleware(handleRoot).ServeHTTP(rr, req)
	case "/time":
		loggingMiddleware(handleTime).ServeHTTP(rr, req)
	case "/echo":
		if method == http.MethodGet {
			loggingMiddleware(handleEchoQuery).ServeHTTP(rr, req)
		} else if method == http.MethodPost {
			loggingMiddleware(handleEchoPost).ServeHTTP(rr, req)
		}
	case "/headers":
		loggingMiddleware(handleHeaders).ServeHTTP(rr, req)
	case "/health":
		loggingMiddleware(handleHealth).ServeHTTP(rr, req)
	default:
		t.Fatalf("Unknown path: %s", path)
	}

	return rr
}

func TestEndpoints(t *testing.T) {
	tests := []struct {
		method       string
		path         string
		expectString string
		body         []byte
	}{
		{"GET", "/", "Welcome to the Go Web Service!", nil},
		{"GET", "/time", `"time":`, nil},
		{"GET", "/echo", `"echo":"No message provided"`, nil},
		{"POST", "/echo", `"echo":{"msg":"hello"}`, []byte(`{"msg":"hello"}`)},
		{"GET", "/headers", `"User-Agent"`, nil},
		{"GET", "/health", "OK", nil},
	}

	for _, tt := range tests {
		rr := executeRequest(t, tt.method, tt.path, tt.body)

		if rr.Code != http.StatusOK {
			t.Errorf("%s %s: expected status 200, got %d", tt.method, tt.path, rr.Code)
		}

		if !strings.Contains(rr.Body.String(), tt.expectString) {
			t.Errorf("%s %s: expected body to contain %q, got %q", tt.method, tt.path, tt.expectString, rr.Body.String())
		}
	}
}
