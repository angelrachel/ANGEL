package assessmentapi

import (
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
	for _, path := range []string{"/api/v1/checks", "/api/v1/capabilities"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, path, nil))
		if response.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s status = %d", path, response.Code)
		}
	}
}
