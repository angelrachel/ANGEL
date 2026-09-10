package execution

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistryRoutesHTTPChecksToHTTPAdapter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	registry := Registry{adapters: []Adapter{HTTPMetadataAdapter{Client: server.Client()}, FixtureAdapter{}}}
	output, err := registry.Run(context.Background(), Input{JobID: "job-http", Target: server.URL, CheckID: "web-security-headers"})
	if err != nil {
		t.Fatal(err)
	}
	if output.Status != "COMPLETED" || output.Result.CheckID != "web-security-headers" {
		t.Fatalf("output=%#v", output)
	}
}

func TestRegistryRejectsUnscopedTargetScheme(t *testing.T) {
	_, err := NewRegistry().Run(context.Background(), Input{Target: "file:///tmp/a", CheckID: "web-security-headers"})
	if err == nil {
		t.Fatal("unsupported target scheme accepted")
	}
}
