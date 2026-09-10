package persistence

import (
	"fmt"
	"sort"
	"strings"

	"ANGEL/src/assessment/evidence"
	"ANGEL/src/assessment/governance"
)

type RestorePolicy struct {
	ReplaceExisting       bool
	RequireAuditChain     bool
	RequireEvidenceHashes bool
}
type RestoreResult struct {
	Engagements int    `json:"engagements"`
	Jobs        int    `json:"jobs"`
	Findings    int    `json:"findings"`
	Evidence    int    `json:"evidence"`
	AuditEvents int    `json:"audit_events"`
	Digest      string `json:"digest"`
}

func Validate(snapshot Snapshot, publicKey []byte, policy RestorePolicy) (RestoreResult, error) {
	if len(publicKey) != 32 {
		return RestoreResult{}, fmt.Errorf("snapshot public key is invalid")
	}
	if err := Verify(snapshot, publicKey); err != nil {
		return RestoreResult{}, err
	}
	if !policy.ReplaceExisting && snapshot.Version < 1 {
		return RestoreResult{}, fmt.Errorf("snapshot version is not restorable")
	}
	if policy.RequireAuditChain && !governance.Verify(snapshot.Audit) {
		return RestoreResult{}, fmt.Errorf("audit chain verification failed")
	}
	if policy.RequireEvidenceHashes {
		for _, item := range snapshot.Evidence {
			if strings.TrimSpace(item.SHA256) == "" || len(item.SHA256) != 64 {
				return RestoreResult{}, fmt.Errorf("evidence hash is incomplete")
			}
		}
	}
	if duplicateEngagement(snapshot) || duplicateJobs(snapshot) || duplicateFindings(snapshot) || duplicateEvidence(snapshot) {
		return RestoreResult{}, fmt.Errorf("snapshot contains duplicate identifiers")
	}
	return RestoreResult{Engagements: len(snapshot.Engagements), Jobs: len(snapshot.Jobs), Findings: len(snapshot.Findings), Evidence: len(snapshot.Evidence), AuditEvents: len(snapshot.Audit), Digest: snapshot.Digest}, nil
}
func duplicateEngagement(s Snapshot) bool {
	ids := map[string]bool{}
	for _, item := range s.Engagements {
		if item.ID == "" || ids[item.ID] {
			return true
		}
		ids[item.ID] = true
	}
	return false
}
func duplicateJobs(s Snapshot) bool {
	ids := map[string]bool{}
	for _, item := range s.Jobs {
		if item.ID == "" || ids[item.ID] {
			return true
		}
		ids[item.ID] = true
	}
	return false
}
func duplicateFindings(s Snapshot) bool {
	ids := map[string]bool{}
	for _, item := range s.Findings {
		if item.ID == "" || ids[item.ID] {
			return true
		}
		ids[item.ID] = true
	}
	return false
}
func duplicateEvidence(s Snapshot) bool {
	ids := map[string]bool{}
	for _, item := range s.Evidence {
		if item.ID == "" || ids[item.ID] {
			return true
		}
		ids[item.ID] = true
	}
	return false
}
func StableIDs(snapshot Snapshot) []string {
	ids := []string{}
	for _, item := range snapshot.Engagements {
		ids = append(ids, item.ID)
	}
	for _, item := range snapshot.Jobs {
		ids = append(ids, item.ID)
	}
	for _, item := range snapshot.Findings {
		ids = append(ids, item.ID)
	}
	for _, item := range snapshot.Evidence {
		ids = append(ids, item.ID)
	}
	sort.Strings(ids)
	return ids
}

var _ = evidence.Bundle{}
