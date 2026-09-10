package risk

import (
	"testing"
	"time"
)

func TestRemediationSupportsPartialAndRiskAcceptedStates(t *testing.T) {
	tracker := NewTracker()
	now := time.Now().UTC()
	item, err := tracker.Create("finding-2", "owner@example.test", "Apply staged controls", now.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	for _, status := range []string{StatusInProgress, StatusPartiallyRemediated, StatusReadyForRetest} {
		if err := tracker.Update(item.ID, status, now); err != nil {
			t.Fatalf("status %s: %v", status, err)
		}
	}
	if err := tracker.Update(item.ID, StatusInProgress, now); err != nil {
		t.Fatal(err)
	}
	if err := tracker.Update(item.ID, StatusRiskAccepted, now); err != nil {
		t.Fatal(err)
	}
	current, _ := tracker.Get(item.ID)
	if current.Status != StatusRiskAccepted {
		t.Fatalf("status=%s", current.Status)
	}
	if err := tracker.Update(item.ID, StatusOpen, now); err == nil {
		t.Fatal("risk accepted remediation was mutable")
	}
}
