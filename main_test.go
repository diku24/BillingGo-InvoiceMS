package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLiveHealthEndpoint(t *testing.T) {
	router := NewRouter()
	request := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != `{"status":"ok"}` {
		t.Fatalf("body = %q, want %q", response.Body.String(), `{"status":"ok"}`)
	}
}

func TestPingUsesConfiguredPort(t *testing.T) {
	router := NewRouter(":9393")
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != `"Server is running on Port :9393"` {
		t.Fatalf("body = %q, want %q", response.Body.String(), `"Server is running on Port :9393"`)
	}
}
