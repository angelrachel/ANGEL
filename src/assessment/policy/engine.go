package policy

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"ANGEL/src/assessment/domain"
)

type Decision struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
}

type Engine struct {
	mu          sync.RWMutex
	engagements map[string]domain.Engagement
	scopes      map[string][]domain.ScopeEntry
	budgets     map[string]int
	used        map[string]int
}

func NewEngine() *Engine {
	return &Engine{engagements: make(map[string]domain.Engagement), scopes: make(map[string][]domain.ScopeEntry), budgets: make(map[string]int), used: make(map[string]int)}
}

func (e *Engine) RegisterEngagement(eng domain.Engagement, scope []domain.ScopeEntry, budget int) error {
	if strings.TrimSpace(eng.ID) == "" || strings.TrimSpace(eng.Organization) == "" {
		return fmt.Errorf("engagement id and organization are required")
	}
	if !eng.Authorized || eng.EndsAt.IsZero() || !eng.EndsAt.After(eng.StartsAt) {
		return fmt.Errorf("engagement is not authorized or has invalid window")
	}
	if budget < 1 {
		return fmt.Errorf("request budget must be positive")
	}
	for _, item := range scope {
		if err := validateScope(item); err != nil {
			return err
		}
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.engagements[eng.ID], e.scopes[eng.ID], e.budgets[eng.ID] = eng, append([]domain.ScopeEntry(nil), scope...), budget
	return nil
}

func (e *Engine) Authorize(engagementID, target, action string, now time.Time) Decision {
	e.mu.Lock()
	defer e.mu.Unlock()
	eng, ok := e.engagements[engagementID]
	if !ok {
		return Decision{Reason: "engagement not found"}
	}
	if eng.Stopped {
		return Decision{Reason: "engagement emergency stop is active"}
	}
	if now.Before(eng.StartsAt) || now.After(eng.EndsAt) {
		return Decision{Reason: "outside engagement window"}
	}
	if e.used[engagementID] >= e.budgets[engagementID] {
		return Decision{Reason: "request budget exhausted"}
	}
	if !allowedTarget(target, e.scopes[engagementID], action) {
		return Decision{Reason: "target or action is outside scope"}
	}
	e.used[engagementID]++
	return Decision{Allowed: true, Reason: "scope, policy, and budget checks passed"}
}

func (e *Engine) Stop(engagementID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	eng := e.engagements[engagementID]
	eng.Stopped = true
	e.engagements[engagementID] = eng
}

func validateScope(item domain.ScopeEntry) error {
	if strings.TrimSpace(item.Value) == "" || strings.TrimSpace(item.Kind) == "" {
		return fmt.Errorf("scope kind and value are required")
	}
	if item.Kind == "url" {
		u, err := url.Parse(item.Value)
		if err != nil || u.Hostname() == "" || (u.Scheme != "https" && u.Scheme != "http") {
			return fmt.Errorf("invalid scoped url")
		}
	}
	if item.Kind == "cidr" {
		if _, _, err := net.ParseCIDR(item.Value); err != nil {
			return fmt.Errorf("invalid scoped cidr")
		}
	}
	if item.Excluded {
		return nil
	}
	if len(item.Actions) == 0 {
		return fmt.Errorf("scope must declare allowed actions")
	}
	return nil
}

func allowedTarget(target string, scopes []domain.ScopeEntry, action string) bool {
	target = strings.TrimSpace(strings.ToLower(target))
	if target == "" || action == "" {
		return false
	}
	for _, item := range scopes {
		if item.Excluded {
			continue
		}
		allowedAction := false
		for _, candidate := range item.Actions {
			if candidate == action {
				allowedAction = true
				break
			}
		}
		if !allowedAction {
			continue
		}
		value := strings.ToLower(strings.TrimSpace(item.Value))
		switch item.Kind {
		case "hostname":
			if target == value || (strings.HasPrefix(value, "*.") && strings.HasSuffix(target, strings.TrimPrefix(value, "*"))) {
				return true
			}
		case "url":
			if u, err := url.Parse(target); err == nil && strings.EqualFold(u.Hostname(), strings.TrimPrefix(value, "https://")) {
				return true
			}
		case "cidr":
			if ip := net.ParseIP(target); ip != nil {
				_, block, _ := net.ParseCIDR(value)
				if block != nil && block.Contains(ip) {
					return true
				}
			}
		case "fixture":
			if strings.HasPrefix(target, "fixture:") && target == value {
				return true
			}
		}
	}
	return false
}
