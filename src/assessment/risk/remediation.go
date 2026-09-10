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

const (
	StatusOpen                = "OPEN"
	StatusInProgress          = "IN_PROGRESS"
	StatusPartiallyRemediated = "PARTIALLY_REMEDIATED"
	StatusReadyForRetest      = "READY_FOR_RETEST"
	StatusRiskAccepted        = "RISK_ACCEPTED"
	StatusClosed              = "CLOSED"
)

func NewTracker() *Tracker {
	return &Tracker{remediation: map[string]Remediation{}, retests: map[string]Retest{}}
}
func (t *Tracker) Create(findingID, owner, plan string, due time.Time) (Remediation, error) {
	if strings.TrimSpace(findingID) == "" || strings.TrimSpace(owner) == "" || strings.TrimSpace(plan) == "" {
		return Remediation{}, fmt.Errorf("finding, owner, and plan are required")
	}
	if due.IsZero() {
		return Remediation{}, fmt.Errorf("due date is required")
	}
	r := Remediation{ID: fmt.Sprintf("rem-%d", time.Now().UnixNano()), FindingID: findingID, Owner: owner, Plan: plan, DueAt: due.UTC(), Status: "OPEN", UpdatedAt: time.Now().UTC()}
	t.mu.Lock()
	t.remediation[r.ID] = r
	t.mu.Unlock()
	return r, nil
}
func (t *Tracker) Update(id, status string, at time.Time) error {
	if status != StatusOpen && status != StatusInProgress && status != StatusPartiallyRemediated && status != StatusReadyForRetest && status != StatusRiskAccepted && status != StatusClosed {
		return fmt.Errorf("invalid remediation status")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	r, ok := t.remediation[id]
	if !ok {
		return fmt.Errorf("remediation not found")
	}
	if r.Status == StatusClosed || r.Status == StatusRiskAccepted {
		return fmt.Errorf("closed remediation cannot be changed")
	}
	if !legalTransition(r.Status, status) {
		return fmt.Errorf("invalid remediation transition %s -> %s", r.Status, status)
	}
	r.Status = status
	r.UpdatedAt = at.UTC()
	t.remediation[id] = r
	return nil
}
func legalTransition(from, to string) bool {
	switch from {
	case StatusOpen:
		return to == StatusInProgress || to == StatusRiskAccepted || to == StatusClosed
	case StatusInProgress:
		return to == StatusPartiallyRemediated || to == StatusReadyForRetest || to == StatusRiskAccepted || to == StatusOpen
	case StatusPartiallyRemediated:
		return to == StatusInProgress || to == StatusReadyForRetest || to == StatusRiskAccepted
	case StatusReadyForRetest:
		return to == StatusInProgress || to == StatusClosed
	default:
		return false
	}
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
func (t *Tracker) Retests(remediationID string) []Retest {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := []Retest{}
	for _, item := range t.retests {
		if item.RemediationID == remediationID {
			out = append(out, item)
		}
	}
	return out
}
