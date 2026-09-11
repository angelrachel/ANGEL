package agent

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Registration identifies a disposable lab agent. It cannot represent a
// production endpoint and only accepts fixture targets.
type Registration struct {
	ID         string    `json:"id"`
	Engagement string    `json:"engagement_id"`
	Target     string    `json:"target"`
	Platform   string    `json:"platform"`
	CreatedAt  time.Time `json:"created_at"`
	Simulation bool      `json:"simulation"`
}

// Task is a bounded, typed simulation task. Arbitrary commands are not part of
// this protocol by design.
type Task struct {
	ID           string            `json:"id"`
	AgentID      string            `json:"agent_id"`
	EngagementID string            `json:"engagement_id"`
	Action       string            `json:"action"`
	FixtureRef   string            `json:"fixture_ref"`
	Parameters   map[string]string `json:"parameters,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
}

// Result is evidence-oriented and always marked as simulation.
type Result struct {
	TaskID      string    `json:"task_id"`
	AgentID     string    `json:"agent_id"`
	Status      string    `json:"status"`
	Simulation  bool      `json:"simulation"`
	Authorized  bool      `json:"authorized"`
	EvidenceRef string    `json:"evidence_ref"`
	Digest      string    `json:"digest"`
	CompletedAt time.Time `json:"completed_at"`
}

type Simulator struct {
	now func() time.Time
}

func NewSimulator() Simulator { return Simulator{now: func() time.Time { return time.Now().UTC() }} }

func (s Simulator) Register(id, engagement, target, platform string) (Registration, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(engagement) == "" || strings.TrimSpace(platform) == "" {
		return Registration{}, errors.New("agent id, engagement, and platform are required")
	}
	if !isFixture(target) {
		return Registration{}, errors.New("lab agent target must use fixture://")
	}
	return Registration{ID: id, Engagement: engagement, Target: target, Platform: platform, CreatedAt: s.clock(), Simulation: true}, nil
}

func (s Simulator) Execute(reg Registration, task Task) (Result, error) {
	if !reg.Simulation || !isFixture(reg.Target) {
		return Result{}, errors.New("registration is not a valid simulation agent")
	}
	if strings.TrimSpace(task.ID) == "" || task.AgentID != reg.ID || task.EngagementID != reg.Engagement {
		return Result{}, errors.New("task does not match simulation agent")
	}
	if !isFixture(task.FixtureRef) {
		return Result{}, errors.New("task fixture reference must use fixture://")
	}
	if !allowedAction(task.Action) {
		return Result{}, fmt.Errorf("action %q is not allowed for lab agent simulation", task.Action)
	}
	digest := sha256.Sum256([]byte(task.ID + "\x00" + task.Action + "\x00" + task.FixtureRef))
	return Result{TaskID: task.ID, AgentID: reg.ID, Status: "completed", Simulation: true, Authorized: true, EvidenceRef: task.FixtureRef, Digest: hex.EncodeToString(digest[:]), CompletedAt: s.clock()}, nil
}

func (s Simulator) clock() time.Time {
	if s.now == nil {
		return time.Now().UTC()
	}
	return s.now().UTC()
}

func isFixture(value string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(value)), "fixture://")
}

func allowedAction(action string) bool {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "observe", "simulate", "collect-evidence", "validate-control":
		return true
	default:
		return false
	}
}
