package orchestrator

import (
"context"
"time"
)

type AttackPath struct {
StartNode string
EndNode   string
Nodes     []string
}

func NewAttackPath(start, end string) *AttackPath {
return &AttackPath{
StartNode: start,
EndNode:   end,
Nodes:     []string{},
}
}

func (a *AttackPath) Generate(ctx context.Context) []string {
a.Nodes = append(a.Nodes, a.StartNode)
a.Nodes = append(a.Nodes, "Recon")
a.Nodes = append(a.Nodes, "Exploit")
a.Nodes = append(a.Nodes, "Pivot")
a.Nodes = append(a.Nodes, a.EndNode)
return a.Nodes
}

func (a *AttackPath) GetPath() []string {
return a.Nodes
}

func (a *AttackPath) GetLastSeen() time.Time {
return time.Now()
}
