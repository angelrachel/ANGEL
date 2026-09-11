package assessmentapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ANGEL/src/assessment/control"
	"ANGEL/src/assessment/policy"
)

func TestRoutesExposeSafeCatalogs(t *testing.T) {
	controller, err := control.NewController(policy.NewEngine())
	if err != nil {
		t.Fatal(err)
	}
	handler := (Handler{Controller: controller}).Routes()
	checks := httptest.NewRecorder()
	handler.ServeHTTP(checks, httptest.NewRequest(http.MethodGet, "/api/v1/checks", nil))
	if checks.Code != http.StatusOK {
		t.Fatalf("checks status = %d", checks.Code)
	}
	var ids []string
	if err := json.Unmarshal(checks.Body.Bytes(), &ids); err != nil || len(ids) < 20 {
		t.Fatalf("unexpected checks response: %v %v", err, checks.Body.String())
	}

	capabilities := httptest.NewRecorder()
	handler.ServeHTTP(capabilities, httptest.NewRequest(http.MethodGet, "/api/v1/capabilities", nil))
	if capabilities.Code != http.StatusOK {
		t.Fatalf("capabilities status = %d", capabilities.Code)
	}
	var policy struct {
		Default string   `json:"default_decision"`
		Allowed []string `json:"allowed"`
		Denied  []string `json:"denied"`
	}
	if err := json.Unmarshal(capabilities.Body.Bytes(), &policy); err != nil {
		t.Fatal(err)
	}
	if policy.Default != "deny" || len(policy.Allowed) == 0 || len(policy.Denied) == 0 {
		t.Fatalf("unexpected capability policy: %+v", policy)
	}
	for _, denied := range policy.Denied {
		if denied == "arbitrary-command" {
			return
		}
	}
	t.Fatal("arbitrary-command must remain denied")
}

func TestCatalogEndpointsRejectWrites(t *testing.T) {
	controller, err := control.NewController(policy.NewEngine())
	if err != nil {
		t.Fatal(err)
	}
	handler := (Handler{Controller: controller}).Routes()
	for _, path := range []string{"/api/v1/checks", "/api/v1/check-catalog", "/api/v1/capabilities"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, path, nil))
		if response.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s status = %d", path, response.Code)
		}
	}
}

func TestTypedCheckCatalogEndpoint(t *testing.T) {
	controller, err := control.NewController(policy.NewEngine())
	if err != nil {
		t.Fatal(err)
	}
	handler := (Handler{Controller: controller}).Routes()
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/check-catalog", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("catalog status = %d", response.Code)
	}
	var catalog []struct {
		ID    string `json:"id"`
		Layer int    `json:"layer"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &catalog); err != nil {
		t.Fatal(err)
	}
	if len(catalog) < 20 || catalog[0].ID == "" || catalog[0].Layer < 1 {
		t.Fatalf("unexpected catalog: %+v", catalog)
	}
}

func TestFixtureExecutionEndpointRunsRegisteredCheck(t *testing.T) {
	controller, err := control.NewController(policy.NewEngine())
	if err != nil {
		t.Fatal(err)
	}
	handler := (Handler{Controller: controller}).Routes()
	body := bytes.NewBufferString(`{"job_id":"job-1","target":"fixture://http/security","check_id":"web-security-headers","fixture":{"headers":{"Content-Security-Policy":"default-src 'self'","X-Content-Type-Options":"nosniff","Referrer-Policy":"no-referrer","Strict-Transport-Security":"max-age=31536000"}}}`)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/executions", body))
	if response.Code != http.StatusCreated {
		t.Fatalf("execution status = %d body=%s", response.Code, response.Body.String())
	}
	var output struct {
		Status  string `json:"status"`
		CheckID string `json:"check_id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatal(err)
	}
	if output.Status != "COMPLETED" || output.CheckID != "web-security-headers" {
		t.Fatalf("unexpected execution output: %+v", output)
	}
}

func TestFixtureExecutionEndpointRejectsProductionTarget(t *testing.T) {
	controller, err := control.NewController(policy.NewEngine())
	if err != nil {
		t.Fatal(err)
	}
	handler := (Handler{Controller: controller}).Routes()
	body := bytes.NewBufferString(`{"job_id":"job-1","target":"https://production.example","check_id":"web-security-headers"}`)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/executions", body))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("execution status = %d", response.Code)
	}
}
