package orchestrator

import (
"context"
"time"
)

type MCPServer struct {
Name     string
Endpoint string
Type     string
}

func NewMCPServer(name, endpoint string) *MCPServer {
return &MCPServer{
Name:     name,
Endpoint: endpoint,
Type:     "MCP",
}
}

func (m *MCPServer) Connect(ctx context.Context) bool {
return true
}

func (m *MCPServer) SendCommand(command string) string {
return "Command sent to " + m.Name + ": " + command
}

func (m *MCPServer) GetStatus() string {
return "MCP Server " + m.Name + " is active"
}

func (m *MCPServer) GetEndpoint() string {
return m.Endpoint
}

func (m *MCPServer) GetType() string {
return m.Type
}

func (m *MCPServer) GetLastSeen() time.Time {
return time.Now()
}

func (m *MCPServer) Retire() {
m.Endpoint = ""
}
