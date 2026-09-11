package policy

import "testing"

func TestCapabilityPolicyIsDenyByDefaultAndSorted(t *testing.T) {
	allowed, denied := CapabilityPolicy()
	if len(allowed) == 0 || len(denied) == 0 {
		t.Fatal("capability policy must not be empty")
	}
	if !CapabilityAllowed("lab-agent-simulation") {
		t.Fatal("lab-agent-simulation should be allowed")
	}
	if !CapabilityDenied("arbitrary-command") || !CapabilityDenied("data-exfiltration") {
		t.Fatal("dangerous capabilities must remain denied")
	}
	for i := 1; i < len(allowed); i++ {
		if allowed[i-1] > allowed[i] {
			t.Fatal("allowed capabilities are not sorted")
		}
	}
	for i := 1; i < len(denied); i++ {
		if denied[i-1] > denied[i] {
			t.Fatal("denied capabilities are not sorted")
		}
	}
}

func TestCapabilityPolicyReturnsCopies(t *testing.T) {
	allowed, denied := CapabilityPolicy()
	allowed[0] = "mutated"
	denied[0] = "mutated"
	freshAllowed, freshDenied := CapabilityPolicy()
	if freshAllowed[0] == "mutated" || freshDenied[0] == "mutated" {
		t.Fatal("capability policy leaked mutable backing arrays")
	}
}
