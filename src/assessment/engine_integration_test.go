package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestControlPlaneHTTPIntegration(t *testing.T) {
	t.Setenv("ANGEL_OPERATOR_KEY", "integration-secret")
	t.Setenv("ANGEL_HOST", "127.0.0.1")
	t.Setenv("ANGEL_PORT", "18081")
	engine := NewEngine()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = engine.Start(ctx) }()
	client := &http.Client{Timeout: 500 * time.Millisecond}
	base := "http://127.0.0.1:18081"
	deadline := time.Now().Add(3 * time.Second)
	for {
		if time.Now().After(deadline) {
			t.Fatal("control plane did not start")
		}
		response, err := client.Get(base + "/healthz")
		if err == nil {
			response.Body.Close()
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	request, _ := http.NewRequest(http.MethodGet, base+"/api/v1/checks", nil)
	request.Header.Set("Authorization", "Bearer integration-secret")
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("checks status=%d", response.StatusCode)
	}
	var checks []string
	if err := json.NewDecoder(response.Body).Decode(&checks); err != nil || len(checks) < 12 {
		t.Fatalf("checks=%d err=%v", len(checks), err)
	}
	payload := `{"check_id":"web-security-headers","input":{"id":"integration","target":"fixture://lab/web","headers":{}}}`
	request, _ = http.NewRequest(http.MethodPost, base+"/api/v1/checks", strings.NewReader(payload))
	request.Header.Set("Authorization", "Bearer integration-secret")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("check evaluation status=%d", response.StatusCode)
	}
}
