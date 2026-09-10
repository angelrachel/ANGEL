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
	engagementPayload := `{"engagement":{"id":"integration-eng","name":"Integration","organization_id":"org","authorized":true,"starts_at":"2020-01-01T00:00:00Z","ends_at":"2099-01-01T00:00:00Z","policy_hash":"p"},"scope":[{"id":"scope-1","engagement_id":"integration-eng","kind":"fixture","value":"fixture://lab/web","actions":["surface-map"]}],"budget":10}`
	response = doJSON(t, client, base+"/api/v1/engagements", http.MethodPost, engagementPayload, "Bearer integration-secret")
	if response.StatusCode != http.StatusCreated {
		response.Body.Close()
		t.Fatalf("engagement status=%d", response.StatusCode)
	}
	response.Body.Close()
	jobPayload := `{"engagement_id":"integration-eng","task_type":"surface-map","target":"fixture://lab/web","action":"surface-map"}`
	response = doJSON(t, client, base+"/api/v1/assessment-jobs", http.MethodPost, jobPayload, "Bearer integration-secret")
	if response.StatusCode != http.StatusAccepted {
		response.Body.Close()
		t.Fatalf("job status=%d", response.StatusCode)
	}
	var job struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&job); err != nil {
		response.Body.Close()
		t.Fatal(err)
	}
	response.Body.Close()
	executionPayload := `{"job_id":"` + job.ID + `","target":"fixture://lab/web","check_id":"web-security-headers","fixture":{"body":"integration"}}`
	response = doJSON(t, client, base+"/api/v1/executions", http.MethodPost, executionPayload, "Bearer integration-secret")
	if response.StatusCode != http.StatusCreated {
		response.Body.Close()
		t.Fatalf("execution status=%d", response.StatusCode)
	}
	response.Body.Close()
	request, _ = http.NewRequest(http.MethodGet, base+"/api/v1/audit/verify", nil)
	request.Header.Set("Authorization", "Bearer integration-secret")
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		t.Fatalf("audit verify status=%d", response.StatusCode)
	}
	response.Body.Close()
}

func doJSON(t *testing.T, client *http.Client, target, method, payload, auth string) *http.Response {
	t.Helper()
	request, err := http.NewRequest(method, target, strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", auth)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	return response
}
