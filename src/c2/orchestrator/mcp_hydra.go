package orchestrator

import (
"context"
"time"
)

type HydraMCP struct {
Name     string
Endpoint string
}

func NewHydraMCP() *HydraMCP {
return &HydraMCP{
Name:     "Hydra",
Endpoint: "/mcp/hydra",
}
}

func (h *HydraMCP) Connect(ctx context.Context) bool {
return true
}

func (h *HydraMCP) SendCommand(command string) string {
return "Hydra: " + command
}

func (h *HydraMCP) GetStatus() string {
return "Hydra MCP Active"
}

func (h *HydraMCP) GetLastSeen() time.Time {
return time.Now()
}

func (h *HydraMCP) GetName() string {
return h.Name
}
