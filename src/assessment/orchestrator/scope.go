package orchestrator

import (
	"fmt"
	"strings"
	"time"
)

// ActionClass describes the operational risk boundary of a requested action.
// The policy layer does not execute actions; it only decides whether a caller
// may submit one to a downstream simulator or assessment module.
type ActionClass string

const (
	ActionPassive    ActionClass = "passive"
	ActionSimulation ActionClass = "simulation"
	ActionActive     ActionClass = "active"
)

type EngagementScope struct {
	ID                string
	Authorized        bool
	StartsAt          time.Time
	EndsAt            time.Time
	Targets           []string
	AllowedTechniques []string
	RulesOfEngagement string
	EmergencyContact  string
}

type ActionRequest struct {
	Target      string
	Technique   string
	Class       ActionClass
	ApprovalID  string
	RequestedBy string
}

type AuthorizationDecision struct {
	Allowed bool
	Reason  string
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func contains(values []string, value string) bool {
	value = normalize(value)
	for _, candidate := range values {
		if normalize(candidate) == value {
			return true
		}
	}
	return false
}

// Authorize evaluates one action against a signed/recorded engagement scope.
// It intentionally fails closed for missing scope data, expired scopes,
// unknown targets, unknown techniques, and active actions without approval.
func (s EngagementScope) Authorize(req ActionRequest, now time.Time) AuthorizationDecision {
	if !s.Authorized {
		return AuthorizationDecision{Reason: "engagement is not authorized"}
	}
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.RulesOfEngagement) == "" {
		return AuthorizationDecision{Reason: "scope identity or Rules of Engagement is missing"}
	}
	if s.StartsAt.IsZero() || s.EndsAt.IsZero() || !now.Before(s.EndsAt) || now.Before(s.StartsAt) {
		return AuthorizationDecision{Reason: "engagement scope is outside its validity window"}
	}
	if strings.TrimSpace(req.RequestedBy) == "" {
		return AuthorizationDecision{Reason: "requester identity is missing"}
	}
	if !contains(s.Targets, req.Target) {
		return AuthorizationDecision{Reason: fmt.Sprintf("target %q is outside the engagement scope", req.Target)}
	}
	if !contains(s.AllowedTechniques, req.Technique) {
		return AuthorizationDecision{Reason: fmt.Sprintf("technique %q is not allowed by the Rules of Engagement", req.Technique)}
	}
	if req.Class != ActionPassive && req.Class != ActionSimulation && req.Class != ActionActive {
		return AuthorizationDecision{Reason: "unknown action class"}
	}
	if req.Class == ActionActive && strings.TrimSpace(req.ApprovalID) == "" {
		return AuthorizationDecision{Reason: "active action requires an explicit approval ID"}
	}
	return AuthorizationDecision{Allowed: true, Reason: "action is within the authorized engagement scope"}
}
