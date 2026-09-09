package orchestrator

import (
	"context"
	"time"
)

type LangGraph struct {
	Nodes     []string
	Edges     map[string][]string
	Current   string
	StartNode string
	EndNodes  []string
}

func NewLangGraph() *LangGraph {
	return &LangGraph{
		Nodes:     []string{},
		Edges:     make(map[string][]string),
		StartNode: "",
		EndNodes:  []string{},
	}
}

func (l *LangGraph) AddNode(node string) {
	l.Nodes = append(l.Nodes, node)
}

func (l *LangGraph) AddEdge(from, to string) {
	l.Edges[from] = append(l.Edges[from], to)
}

func (l *LangGraph) SetStartNode(node string) {
	l.StartNode = node
	l.Current = node
}

func (l *LangGraph) SetEndNodes(nodes []string) {
	l.EndNodes = nodes
}

func (l *LangGraph) Traverse(ctx context.Context) []string {
	var path []string
	current := l.StartNode
	for current != "" {
		path = append(path, current)
		nextNodes := l.Edges[current]
		if len(nextNodes) == 0 {
			break
		}
		current = nextNodes[0]
	}
	return path
}

func (l *LangGraph) GetCurrentNode() string {
	return l.Current
}

func (l *LangGraph) GetNodes() []string {
	return l.Nodes
}

func (l *LangGraph) GetEdges() map[string][]string {
	return l.Edges
}

func (l *LangGraph) GetLastSeen() time.Time {
	return time.Now()
}
