package orchestrator

import "sync"

type Agent struct {
Name    string
Type    string
Handler func(task Task) Result
}

type AgentRegistry struct {
mu     sync.RWMutex
agents map[string]*Agent
}

func NewAgentRegistry() *AgentRegistry {
return &AgentRegistry{agents: make(map[string]*Agent)}
}

func (a *AgentRegistry) Register(agent *Agent) {
a.mu.Lock()
defer a.mu.Unlock()
a.agents[agent.Name] = agent
}

func (a *AgentRegistry) Get(name string) *Agent {
a.mu.RLock()
defer a.mu.RUnlock()
return a.agents[name]
}

func (a *AgentRegistry) List() []*Agent {
a.mu.RLock()
defer a.mu.RUnlock()
var agents []*Agent
for _, agent := range a.agents {
agents = append(agents, agent)
}
return agents
}
