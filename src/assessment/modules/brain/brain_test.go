package brain

import (
	"testing"
	"time"
)

func TestAutonomousDecisionAnalyzeSafeEnvironment(t *testing.T) {
	a := NewAutonomousDecision()
	decision := a.Analyze("safe")
	if decision.Action != "continue_operation" {
		t.Fatalf("expected continue_operation, got %q", decision.Action)
	}
	if decision.Timestamp == "" {
		t.Fatal("expected timestamp")
	}
	parsed, err := time.Parse(time.RFC3339, decision.Timestamp)
	if err != nil {
		t.Fatalf("invalid timestamp: %v", err)
	}
	if parsed.IsZero() {
		t.Fatal("timestamp is zero")
	}
}

func TestAutonomousDecisionAnalyzeUnsafeEnvironment(t *testing.T) {
	a := NewAutonomousDecision()
	decision := a.Analyze("unsafe")
	if decision.Action != "sleep_longer" {
		t.Fatalf("expected sleep_longer, got %q", decision.Action)
	}
}

func TestAutonomousDecisionExecuteRecordsState(t *testing.T) {
	a := NewAutonomousDecision()
	result := a.Execute("test_action")
	if result != "executed: test_action" {
		t.Fatalf("unexpected result: %q", result)
	}
	if a.State["last_action"] != "test_action" {
		t.Fatal("state was not updated")
	}
}

func TestRiskAssessorHighRisk(t *testing.T) {
	assessor := NewRiskAssessor(50)
	result := assessor.Assess(75)
	if result.Level != "high" {
		t.Fatalf("expected high, got %q", result.Level)
	}
	if result.Score != 75 {
		t.Fatalf("expected score 75, got %d", result.Score)
	}
}

func TestRiskAssessorLowRisk(t *testing.T) {
	assessor := NewRiskAssessor(50)
	result := assessor.Assess(25)
	if result.Level != "low" {
		t.Fatalf("expected low, got %q", result.Level)
	}
}

func TestRiskAssessorThresholdBorderline(t *testing.T) {
	assessor := NewRiskAssessor(50)
	result := assessor.Assess(50)
	if result.Level != "high" {
		t.Fatalf("expected high at threshold, got %q", result.Level)
	}
}

func TestRiskAssessorTimestamp(t *testing.T) {
	assessor := NewRiskAssessor(50)
	ts := assessor.GetTimestamp()
	if ts == "" {
		t.Fatal("expected timestamp")
	}
	if _, err := time.Parse(time.RFC3339, ts); err != nil {
		t.Fatalf("invalid timestamp: %v", err)
	}
}
