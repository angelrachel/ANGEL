package domain

import "time"

type Engagement struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Organization string    `json:"organization_id"`
	Authorized   bool      `json:"authorized"`
	StartsAt     time.Time `json:"starts_at"`
	EndsAt       time.Time `json:"ends_at"`
	PolicyHash   string    `json:"policy_hash"`
	Stopped      bool      `json:"stopped"`
}

type ScopeEntry struct {
	ID           string   `json:"id"`
	EngagementID string   `json:"engagement_id"`
	Kind         string   `json:"kind"`
	Value        string   `json:"value"`
	Ports        []int    `json:"ports,omitempty"`
	Actions      []string `json:"actions"`
	Excluded     bool     `json:"excluded"`
}

type Asset struct {
	ID           string    `json:"id"`
	EngagementID string    `json:"engagement_id"`
	Canonical    string    `json:"canonical"`
	Kind         string    `json:"kind"`
	Owner        string    `json:"owner,omitempty"`
	Criticality  string    `json:"criticality"`
	Confidence   float64   `json:"confidence"`
	LastSeen     time.Time `json:"last_seen"`
}

type AssessmentJob struct {
	ID           string    `json:"id"`
	EngagementID string    `json:"engagement_id"`
	TaskType     string    `json:"task_type"`
	Target       string    `json:"target"`
	PolicyHash   string    `json:"policy_hash"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	StartedAt    time.Time `json:"started_at,omitempty"`
	FinishedAt   time.Time `json:"finished_at,omitempty"`
	Signature    string    `json:"signature"`
}

type Observation struct {
	ID         string    `json:"id"`
	JobID      string    `json:"job_id"`
	CheckID    string    `json:"check_id"`
	Target     string    `json:"target"`
	Status     string    `json:"status"`
	Confidence float64   `json:"confidence"`
	Summary    string    `json:"summary"`
	EvidenceID string    `json:"evidence_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Finding struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Severity       string   `json:"severity"`
	Confidence     float64  `json:"confidence"`
	AffectedAssets []string `json:"affected_assets"`
	Impact         string   `json:"impact"`
	EvidenceIDs    []string `json:"evidence_ids"`
	Remediation    string   `json:"remediation"`
	Status         string   `json:"status"`
}

type AuditEvent struct {
	ID        string    `json:"id"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Object    string    `json:"object"`
	ObjectID  string    `json:"object_id"`
	Timestamp time.Time `json:"timestamp"`
	Hash      string    `json:"hash"`
}
