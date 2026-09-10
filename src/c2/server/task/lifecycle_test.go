package task

import (
	"testing"
	"time"
)

func TestLifecycleAllowsApprovedExecutionPath(t *testing.T) {
	now := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	life, err := NewLifecycle("task-1", now)
	if err != nil {
		t.Fatal(err)
	}
	for _, next := range []State{StateApproved, StateQueued, StateRunning, StateCompleted} {
		life, err = life.Transition(next, "approved by policy", now.Add(time.Minute))
		if err != nil {
			t.Fatalf("transition to %s: %v", next, err)
		}
	}
	if life.Version != 5 || life.State != StateCompleted {
		t.Fatalf("unexpected lifecycle: %#v", life)
	}
}

func TestLifecycleRejectsInvalidTransition(t *testing.T) {
	life, err := NewLifecycle("task-1", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := life.Transition(StateRunning, "", time.Now().UTC()); err == nil {
		t.Fatal("expected created-to-running transition to fail")
	}
}

func TestLifecycleRequiresReasonsForTerminalFailureStates(t *testing.T) {
	life, err := NewLifecycle("task-1", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	life, err = life.Transition(StateApproved, "policy approved", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := life.Transition(StateRejected, "", time.Now().UTC()); err == nil {
		t.Fatal("expected rejection reason requirement")
	}
}

func TestLifecycleCompensatesFailedTask(t *testing.T) {
	now := time.Now().UTC()
	life, err := NewLifecycle("task-1", now)
	if err != nil {
		t.Fatal(err)
	}
	for _, next := range []State{StateApproved, StateQueued, StateRunning, StateFailed, StateCompensated} {
		reason := "policy approved"
		if next == StateFailed {
			reason = "worker timeout"
		}
		if next == StateCompensated {
			reason = "rollback verified"
		}
		life, err = life.Transition(next, reason, now)
		if err != nil {
			t.Fatalf("transition to %s: %v", next, err)
		}
	}
	if life.State != StateCompensated || life.Reason != "rollback verified" {
		t.Fatalf("unexpected compensation state: %#v", life)
	}
}
