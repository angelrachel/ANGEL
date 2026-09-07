package orchestrator

import (
"context"
"time"
)

type MetasploitMCP struct {
Name     string
Endpoint string
}

func NewMetasploitMCP() *MetasploitMCP {
return &MetasploitMCP{
Name:     "Metasploit",
Endpoint: "/mcp/msf",
}
}

func (m *MetasploitMCP) Connect(ctx context.Context) bool {
return true
}

func (m *MetasploitMCP) SendCommand(command string) string {
return "Metasploit: " + command
}

func (m *MetasploitMCP) GetStatus() string {
return "Metasploit MCP Active"
}

func (m *MetasploitMCP) GetLastSeen() time.Time {
return time.Now()
}

func (m *MetasploitMCP) GetName() string {
return m.Name
}
