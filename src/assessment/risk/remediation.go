package risk

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Remediation struct {
	ID        string    `json:"id"`
	FindingID string    `json:"finding_id"`
	Owner     string    `json:"owner"`
	Plan      string    `json:"plan"`
	DueAt     time.Time `json:"due_at"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Retest struct {
	ID            string    `json:"id"`
	RemediationID string    `json:"remediation_id"`
	Passed        bool      `json:"passed"`
	EvidenceID    string    `json:"evidence_id"`
	Notes         string    `json:"notes"`
	TestedAt      time.Time `json:"tested_at"`
}
type Tracker struct {
	mu          sync.RWMutex
	remediation map[string]Remediation
	retests     map[string]Retest
}

func NewTracker() *Tracker {
	return &Tracker{remediation: map[string]Remediation{}, retests: map[string]Retest{}}
}
func (t *Tracker) Create(findingID, owner, plan string, due time.Time) (Remediation, error) {
	if strings.TrimSpace(findingID) == "" || strings.TrimSpace(owner) == "" || strings.TrimSpace(plan) == "" {
		return Remediation{}, fmt.Errorf("finding, owner, and plan are required")
	}
	r := Remediation{ID: fmt.Sprintf("rem-%d", time.Now().UnixNano()), FindingID: findingID, Owner: owner, Plan: plan, DueAt: due.UTC(), Status: "OPEN", UpdatedAt: time.Now().UTC()}
	t.mu.Lock()
	t.remediation[r.ID] = r
	t.mu.Unlock()
	return r, nil
}
func (t *Tracker) Update(id, status string, at time.Time) error {
	if status != "OPEN" && status != "IN_PROGRESS" && status != "READY_FOR_RETEST" && status != "CLOSED" {
		return fmt.Errorf("invalid remediation status")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	r, ok := t.remediation[id]
	if !ok {
		return fmt.Errorf("remediation not found")
	}
	if r.Status == "CLOSED" {
		return fmt.Errorf("closed remediation cannot be changed")
	}
	r.Status = status
	r.UpdatedAt = at.UTC()
	t.remediation[id] = r
	return nil
}
func (t *Tracker) RecordRetest(remediationID, evidenceID, notes string, passed bool, at time.Time) (Retest, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	r, ok := t.remediation[remediationID]
	if !ok {
		return Retest{}, fmt.Errorf("remediation not found")
	}
	if r.Status != "READY_FOR_RETEST" {
		return Retest{}, fmt.Errorf("remediation is not ready for retest")
	}
	if strings.TrimSpace(evidenceID) == "" {
		return Retest{}, fmt.Errorf("retest evidence is required")
	}
	result := Retest{ID: fmt.Sprintf("retest-%d", time.Now().UnixNano()), RemediationID: remediationID, Passed: passed, EvidenceID: evidenceID, Notes: notes, TestedAt: at.UTC()}
	if passed {
		r.Status = "CLOSED"
	} else {
		r.Status = "IN_PROGRESS"
	}
	r.UpdatedAt = at.UTC()
	t.remediation[remediationID] = r
	t.retests[result.ID] = result
	return result, nil
}
func (t *Tracker) Get(id string) (Remediation, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.remediation[id]
	return r, ok
}

func (t *Tracker) List() []Remediation {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Remediation, 0, len(t.remediation))
	for _, item := range t.remediation {
		out = append(out, item)
	}
	return out
}
