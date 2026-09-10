package correlation

import (
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Label       string `json:"label"`
	Criticality int    `json:"criticality"`
}
type Edge struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	Relation   string  `json:"relation"`
	Confidence float64 `json:"confidence"`
}
type Path struct {
	Nodes     []string `json:"nodes"`
	Relations []string `json:"relations"`
	Score     float64  `json:"score"`
	Priority  string   `json:"priority"`
}
type Graph struct {
	Nodes map[string]Node `json:"nodes"`
	Edges []Edge          `json:"edges"`
}

func New() *Graph { return &Graph{Nodes: map[string]Node{}, Edges: []Edge{}} }
func (g *Graph) AddNode(node Node) error {
	if strings.TrimSpace(node.ID) == "" || strings.TrimSpace(node.Kind) == "" {
		return fmt.Errorf("node identity is required")
	}
	if node.Criticality < 0 || node.Criticality > 10 {
		return fmt.Errorf("criticality must be between 0 and 10")
	}
	g.Nodes[node.ID] = node
	return nil
}
func (g *Graph) AddEdge(edge Edge) error {
	if _, ok := g.Nodes[edge.From]; !ok {
		return fmt.Errorf("source node not found")
	}
	if _, ok := g.Nodes[edge.To]; !ok {
		return fmt.Errorf("target node not found")
	}
	if edge.Confidence < 0 || edge.Confidence > 1 {
		return fmt.Errorf("edge confidence must be between 0 and 1")
	}
	g.Edges = append(g.Edges, edge)
	return nil
}
func (g *Graph) Paths(source string) []Path {
	if _, ok := g.Nodes[source]; !ok {
		return nil
	}
	adj := map[string][]Edge{}
	for _, edge := range g.Edges {
		adj[edge.From] = append(adj[edge.From], edge)
	}
	var out []Path
	var visit func(string, []string, []string, float64, map[string]bool)
	visit = func(current string, nodes, relations []string, score float64, seen map[string]bool) {
		if len(nodes) > 1 && len(adj[current]) == 0 {
			priority := "HIGH"
			if score >= 8 {
				priority = "P1"
			}
			if score >= 9.2 && g.Nodes[nodes[len(nodes)-1]].Criticality >= 9 {
				priority = "P0"
			}
			out = append(out, Path{Nodes: append([]string(nil), nodes...), Relations: append([]string(nil), relations...), Score: score, Priority: priority})
			return
		}
		for _, edge := range adj[current] {
			if seen[edge.To] {
				continue
			}
			next := map[string]bool{}
			for k, v := range seen {
				next[k] = v
			}
			next[edge.To] = true
			nextScore := score + float64(g.Nodes[edge.To].Criticality)*edge.Confidence
			visit(edge.To, append(nodes, edge.To), append(relations, edge.Relation), nextScore, next)
		}
	}
	visit(source, []string{source}, nil, float64(g.Nodes[source].Criticality), map[string]bool{source: true})
	sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}
