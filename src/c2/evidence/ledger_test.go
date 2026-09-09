package evidence

import "testing"

func TestVerifyChainAcceptsLedger(t *testing.T) {
	ledger := NewLedger()
	ledger.AddRecord(`{"event":"scope-created"}`)
	ledger.AddRecord(`{"event":"assessment-started"}`)

	if !VerifyChain(ledger.Records) {
		t.Fatal("expected generated ledger to verify")
	}
}

func TestVerifyChainRejectsTampering(t *testing.T) {
	ledger := NewLedger()
	ledger.AddRecord(`{"event":"scope-created"}`)
	ledger.AddRecord(`{"event":"assessment-started"}`)
	ledger.Records[1].Data = `{"event":"unauthorized-change"}`

	if VerifyChain(ledger.Records) {
		t.Fatal("expected modified record to fail verification")
	}
}

func TestVerifyChainRejectsBrokenSequenceAndParent(t *testing.T) {
	ledger := NewLedger()
	ledger.AddRecord(`{"event":"scope-created"}`)
	ledger.AddRecord(`{"event":"assessment-started"}`)
	ledger.Records[1].Sequence = 7
	ledger.Records[1].Parent = "wrong-parent"

	if VerifyChain(ledger.Records) {
		t.Fatal("expected broken sequence and parent link to fail verification")
	}
}

func TestParentChildValidation(t *testing.T) {
	valid := ParentChild{Parent: "scope-1", Child: "task-1"}
	if !valid.Validate() {
		t.Fatal("expected distinct non-empty IDs to be valid")
	}

	for _, invalid := range []ParentChild{
		{},
		{Parent: "scope-1"},
		{Child: "task-1"},
		{Parent: "same", Child: "same"},
	} {
		if invalid.Validate() {
			t.Fatalf("expected invalid relation %#v to be rejected", invalid)
		}
	}
}
