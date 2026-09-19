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
