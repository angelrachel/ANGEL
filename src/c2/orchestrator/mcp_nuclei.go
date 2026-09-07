package orchestrator

import (
"context"
"time"
)

type NucleiMCP struct {
Name     string
Endpoint string
}

func NewNucleiMCP() *NucleiMCP {
return &NucleiMCP{
Name:     "Nuclei",
Endpoint: "/mcp/nuclei",
}
}

func (n *NucleiMCP) Connect(ctx context.Context) bool {
return true
}

func (n *NucleiMCP) SendCommand(command string) string {
return "Nuclei scanning: " + command
}

func (n *NucleiMCP) GetStatus() string {
return "Nuclei MCP Active"
}

func (n *NucleiMCP) GetLastSeen() time.Time {
return time.Now()
}

func (n *NucleiMCP) GetName() string {
return n.Name
}
