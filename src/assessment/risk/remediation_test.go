package risk

import (
	"testing"
	"time"
)

func TestRemediationRetestClosesFinding(t *testing.T) {
	tracker := NewTracker()
	now := time.Now().UTC()
	remediation, err := tracker.Create("finding-1", "owner@example.test", "Apply control", now.Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if err := tracker.Update(remediation.ID, "IN_PROGRESS", now); err != nil {
		t.Fatal(err)
	}
	if err := tracker.Update(remediation.ID, "READY_FOR_RETEST", now); err != nil {
		t.Fatal(err)
	}
	retest, err := tracker.RecordRetest(remediation.ID, "evidence-1", "control verified", true, now)
	if err != nil {
		t.Fatal(err)
	}
	if !retest.Passed {
		t.Fatal("retest did not pass")
	}
	current, _ := tracker.Get(remediation.ID)
	if current.Status != "CLOSED" {
		t.Fatalf("status=%s", current.Status)
	}
}
