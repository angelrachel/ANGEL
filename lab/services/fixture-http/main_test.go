package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFixtureHealthEndpoint(t *testing.T) {
	response := httptest.NewRecorder()
	newMux().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if response.Code != http.StatusOK || response.Body.String() == "" {
		t.Fatalf("unexpected health response: %d %q", response.Code, response.Body.String())
	}
}

func TestFixtureSecurityPostureEndpoint(t *testing.T) {
	response := httptest.NewRecorder()
	newMux().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/security-posture", nil))
	for _, header := range []string{"Content-Security-Policy", "X-Content-Type-Options", "Referrer-Policy", "Strict-Transport-Security"} {
		if response.Header().Get(header) == "" {
			t.Errorf("missing %s", header)
		}
	}
}

func TestFixtureSensitiveAndCookieEndpoints(t *testing.T) {
	response := httptest.NewRecorder()
	newMux().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/sensitive", nil))
	if response.Header().Get("Cache-Control") != "private, no-store" {
		t.Fatalf("unexpected cache policy: %q", response.Header().Get("Cache-Control"))
	}

	response = httptest.NewRecorder()
	newMux().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/cookie", nil))
	cookie := response.Header().Get("Set-Cookie")
	for _, marker := range []string{"Secure", "HttpOnly", "SameSite=Lax"} {
		if !contains(cookie, marker) {
			t.Errorf("cookie missing %s: %q", marker, cookie)
		}
	}
}

func TestFixtureCORSIsRestricted(t *testing.T) {
	response := httptest.NewRecorder()
	newMux().ServeHTTP(response, httptest.NewRequest(http.MethodOptions, "/cors", nil))
	if response.Code != http.StatusNoContent || response.Header().Get("Access-Control-Allow-Origin") != "https://fixture.example" {
		t.Fatalf("unexpected CORS response: %d %q", response.Code, response.Header().Get("Access-Control-Allow-Origin"))
	}
}

func contains(value, marker string) bool {
	for i := 0; i+len(marker) <= len(value); i++ {
		if value[i:i+len(marker)] == marker {
			return true
		}
	}
	return false
}
