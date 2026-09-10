package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFixtureHealthEndpoint(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK || recorder.Body.String() == "" {
		t.Fatalf("unexpected health response: %d %q", recorder.Code, recorder.Body.String())
	}
}
