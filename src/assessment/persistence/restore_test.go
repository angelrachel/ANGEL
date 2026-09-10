package persistence

import (
	"testing"
	"time"

	"ANGEL/src/assessment/domain"
	"ANGEL/src/assessment/evidence"
	"ANGEL/src/assessment/governance"
	"ANGEL/src/assessment/reporting"
)

func TestRestoreValidatesSignedSnapshot(t *testing.T) {
	signer, err := NewSigner()
	if err != nil {
		t.Fatal(err)
	}
	audit := governance.NewAuditLog()
	audit.Append("operator", "SNAPSHOT", "snapshot", "s-1", time.Now())
	evidenceSigner, err := evidence.NewSigner()
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := evidence.CreateBundle("eng-1", "job-1", "application/json", []byte(`{"ok":true}`), "0", true, time.Now().UTC(), evidenceSigner)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := Create(Snapshot{Version: 1, Engagements: []domain.Engagement{{ID: "eng-1"}}, Jobs: []domain.AssessmentJob{{ID: "job-1"}}, Findings: []reporting.Finding{{ID: "finding-1"}}, Evidence: []evidence.Bundle{bundle}, Audit: audit.List()}, signer)
	if err != nil {
		t.Fatal(err)
	}
	result, err := Validate(snapshot, signer.Public, RestorePolicy{RequireAuditChain: true, RequireEvidenceHashes: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.Jobs != 1 || len(StableIDs(snapshot)) != 4 {
		t.Fatalf("result=%#v", result)
	}
}
func TestRestoreRejectsDuplicateIDs(t *testing.T) {
	signer, _ := NewSigner()
	snapshot, _ := Create(Snapshot{Version: 1, Engagements: []domain.Engagement{{ID: "duplicate"}, {ID: "duplicate"}}}, signer)
	if _, err := Validate(snapshot, signer.Public, RestorePolicy{}); err == nil {
		t.Fatal("duplicate snapshot accepted")
	}
}
