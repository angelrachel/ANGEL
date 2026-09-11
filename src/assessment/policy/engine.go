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

type SnapshotEntry struct {
	Engagement domain.Engagement   `json:"engagement"`
	Scope      []domain.ScopeEntry `json:"scope"`
	Budget     int                 `json:"budget"`
	Used       int                 `json:"used"`
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

func (e *Engine) Snapshot() []SnapshotEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]SnapshotEntry, 0, len(e.engagements))
	for id, engagement := range e.engagements {
		out = append(out, SnapshotEntry{Engagement: engagement, Scope: append([]domain.ScopeEntry(nil), e.scopes[id]...), Budget: e.budgets[id], Used: e.used[id]})
	}
	return out
}

func (e *Engine) Restore(entries []SnapshotEntry) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, entry := range entries {
		if err := validateEngagement(entry.Engagement, entry.Scope, entry.Budget); err != nil {
			return err
		}
		if entry.Used < 0 || entry.Used > entry.Budget {
			return fmt.Errorf("invalid used budget for engagement %s", entry.Engagement.ID)
		}
		e.engagements[entry.Engagement.ID] = entry.Engagement
		e.scopes[entry.Engagement.ID] = append([]domain.ScopeEntry(nil), entry.Scope...)
		e.budgets[entry.Engagement.ID] = entry.Budget
		e.used[entry.Engagement.ID] = entry.Used
	}
	return nil
}

func (e *Engine) RegisterEngagement(eng domain.Engagement, scope []domain.ScopeEntry, budget int) error {
	if err := validateEngagement(eng, scope, budget); err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.engagements[eng.ID], e.scopes[eng.ID], e.budgets[eng.ID] = eng, append([]domain.ScopeEntry(nil), scope...), budget
	return nil
}

func validateEngagement(eng domain.Engagement, scope []domain.ScopeEntry, budget int) error {
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
	return nil
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

func (e *Engine) Authorize(engagementID, target, action string, now time.Time) Decision {
	if deniedCapability(action) {
		return Decision{Reason: "capability is permanently denied by platform policy"}
	}
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
			if hostnameMatches(target, value) && portMatches(target, item.Ports) {
				return true
			}
		case "url":
			if u, err := url.Parse(target); err == nil {
				scoped, parseErr := url.Parse(value)
				if parseErr == nil && strings.EqualFold(u.Scheme, scoped.Scheme) && strings.EqualFold(u.Hostname(), scoped.Hostname()) && portMatches(target, item.Ports) {
					return true
				}
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
		case "port":
			if portMatches(target, []int{portValue(value)}) {
				return true
			}
		case "path":
			if u, err := url.Parse(target); err == nil && pathMatches(u.Path, value) {
				return true
			}
		}
	}
	return false
}

func hostnameMatches(target, scoped string) bool {
	if target == scoped {
		return true
	}
	if !strings.HasPrefix(scoped, "*.") {
		return false
	}
	suffix := strings.TrimPrefix(scoped, "*.")
	return strings.HasSuffix(target, "."+suffix) && target != suffix
}
func pathMatches(target, scoped string) bool {
	scoped = strings.TrimSpace(scoped)
	if scoped == "" || scoped == "/" {
		return true
	}
	target = "/" + strings.TrimPrefix(target, "/")
	scoped = "/" + strings.TrimPrefix(scoped, "/")
	return target == scoped || strings.HasPrefix(target, strings.TrimSuffix(scoped, "/")+"/")
}
func portMatches(target string, allowed []int) bool {
	if len(allowed) == 0 {
		return true
	}
	u, err := url.Parse(target)
	if err != nil {
		return false
	}
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else if u.Scheme == "http" {
			port = "80"
		}
	}
	actual := portValue(port)
	for _, candidate := range allowed {
		if actual == candidate {
			return true
		}
	}
	return false
}
func portValue(value string) int {
	var port int
	if _, err := fmt.Sscanf(strings.TrimSpace(value), "%d", &port); err != nil {
		return -1
	}
	return port
}
func deniedCapability(action string) bool {
	return CapabilityDenied(strings.ToLower(strings.TrimSpace(action)))
}
