package orchestrator

import (
	"fmt"
	"sync"
)

type Task struct {
	ID      string
	Type    string
	Payload map[string]interface{}
}

type Result struct {
	TaskID string
	Status string
	Output map[string]interface{}
}

type Router struct {
	mu     sync.RWMutex
	agents map[string]interface{}
}

func NewRouter() *Router {
	return &Router{agents: make(map[string]interface{})}
}

func (r *Router) RegisterAgent(name string, agent interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[name] = agent
}

func (r *Router) Dispatch(task Task) Result {
	r.mu.RLock()
	defer r.mu.RUnlock()
	agentName := task.Type
	_, exists := r.agents[agentName]
	if !exists {
		return Result{TaskID: task.ID, Status: "error", Output: map[string]interface{}{"error": "agent not found"}}
	}
	return Result{TaskID: task.ID, Status: "accepted", Output: map[string]interface{}{"agent": agentName, "executed": false, "mode": "observe"}}
}

func (r *Router) GetAgentCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.agents)
}

func FormatRouterStatus(count int) string {
	return fmt.Sprintf("Router aktif dengan %d agen", count)
}
