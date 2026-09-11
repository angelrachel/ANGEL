package agent

import (
	"strings"
	"testing"
	"time"
)

func TestSimulatorRegistersOnlyFixtureTargets(t *testing.T) {
	sim := NewSimulator()
	if _, err := sim.Register("agent-1", "eng-1", "https://production.example", "linux"); err == nil {
		t.Fatal("expected non-fixture registration to fail")
	}
	reg, err := sim.Register("agent-1", "eng-1", "fixture://http/security", "linux")
	if err != nil {
		t.Fatalf("register fixture: %v", err)
	}
	if !reg.Simulation || !strings.HasPrefix(reg.Target, "fixture://") {
		t.Fatalf("unexpected registration: %+v", reg)
	}
}

func TestSimulatorExecutesBoundedActionAndReturnsDigest(t *testing.T) {
	sim := Simulator{now: func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }}
	reg, err := sim.Register("agent-1", "eng-1", "fixture://http/security", "linux")
	if err != nil {
		t.Fatal(err)
	}
	result, err := sim.Execute(reg, Task{ID: "task-1", AgentID: "agent-1", EngagementID: "eng-1", Action: "validate-control", FixtureRef: "fixture://http/security"})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result.Status != "completed" || !result.Simulation || !result.Authorized || result.Digest == "" {
		t.Fatalf("unexpected result: %+v", result)
	}
	if !result.CompletedAt.Equal(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected completion time: %v", result.CompletedAt)
	}
}

func TestSimulatorRejectsDangerousOrMismatchedTasks(t *testing.T) {
	sim := NewSimulator()
	reg, err := sim.Register("agent-1", "eng-1", "fixture://http/security", "linux")
	if err != nil {
		t.Fatal(err)
	}
	cases := []Task{
		{ID: "task-1", AgentID: "agent-1", EngagementID: "eng-1", Action: "arbitrary-command", FixtureRef: "fixture://http/security"},
		{ID: "task-2", AgentID: "agent-2", EngagementID: "eng-1", Action: "observe", FixtureRef: "fixture://http/security"},
		{ID: "task-3", AgentID: "agent-1", EngagementID: "eng-1", Action: "observe", FixtureRef: "https://production.example"},
	}
	for _, task := range cases {
		if _, err := sim.Execute(reg, task); err == nil {
			t.Errorf("expected task rejection: %+v", task)
		}
	}
}
