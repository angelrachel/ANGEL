package orchestrator

import (
	"strings"
	"testing"
	"time"
)

func testScope(now time.Time) EngagementScope {
	return EngagementScope{
		ID:                "eng-2026-001",
		Authorized:        true,
		StartsAt:          now.Add(-time.Hour),
		EndsAt:            now.Add(time.Hour),
		Targets:           []string{"app.lab.example"},
		AllowedTechniques: []string{"http-inventory", "simulated-validation"},
		RulesOfEngagement: "roe-2026-001",
		EmergencyContact:  "soc@example.com",
	}
}

func TestAuthorizeAllowsInScopeSimulation(t *testing.T) {
	now := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	decision := testScope(now).Authorize(ActionRequest{
		Target:      "APP.LAB.EXAMPLE",
		Technique:   "HTTP-INVENTORY",
		Class:       ActionSimulation,
		RequestedBy: "operator@example.com",
	}, now)

	if !decision.Allowed {
		t.Fatalf("expected action to be allowed: %s", decision.Reason)
	}
}

func TestAuthorizeRejectsOutOfScopeAndExpiredRequests(t *testing.T) {
	now := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	cases := []ActionRequest{
		{Target: "outside.example", Technique: "http-inventory", Class: ActionPassive, RequestedBy: "operator"},
		{Target: "app.lab.example", Technique: "unknown", Class: ActionPassive, RequestedBy: "operator"},
	}
	for _, req := range cases {
		if decision := testScope(now).Authorize(req, now); decision.Allowed {
			t.Fatalf("expected request to be rejected: %#v", req)
		}
	}

	expired := testScope(now)
	if decision := expired.Authorize(ActionRequest{
		Target: "app.lab.example", Technique: "http-inventory", Class: ActionPassive, RequestedBy: "operator",
	}, now.Add(2*time.Hour)); decision.Allowed {
		t.Fatal("expected expired scope to be rejected")
	}
}

func TestAuthorizeRequiresApprovalForActiveActions(t *testing.T) {
	now := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	request := ActionRequest{
		Target: "app.lab.example", Technique: "simulated-validation", Class: ActionActive, RequestedBy: "operator",
	}
	decision := testScope(now).Authorize(request, now)
	if decision.Allowed || !strings.Contains(decision.Reason, "approval ID") {
		t.Fatalf("expected approval gate, got %#v", decision)
	}

	request.ApprovalID = "approval-001"
	if decision = testScope(now).Authorize(request, now); !decision.Allowed {
		t.Fatalf("expected approved action to pass: %s", decision.Reason)
	}
}

func TestAuthorizeFailsClosedForMissingAuthorizationData(t *testing.T) {
	now := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	scope := testScope(now)
	scope.Authorized = false
	if decision := scope.Authorize(ActionRequest{
		Target: "app.lab.example", Technique: "http-inventory", Class: ActionPassive, RequestedBy: "operator",
	}, now); decision.Allowed {
		t.Fatal("expected unauthorized scope to be rejected")
	}

	scope = testScope(now)
	scope.RulesOfEngagement = ""
	if decision := scope.Authorize(ActionRequest{
		Target: "app.lab.example", Technique: "http-inventory", Class: ActionPassive, RequestedBy: "operator",
	}, now); decision.Allowed {
		t.Fatal("expected scope without ROE to be rejected")
	}
}
