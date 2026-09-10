package correlation

import "testing"

func TestGraphBuildsPrioritizedPath(t *testing.T) {
	graph := New()
	for _, node := range []Node{{ID: "user", Kind: "identity", Criticality: 9}, {ID: "app", Kind: "application", Criticality: 9}, {ID: "db", Kind: "database", Criticality: 10}} {
		if err := graph.AddNode(node); err != nil {
			t.Fatal(err)
		}
	}
	if err := graph.AddEdge(Edge{From: "user", To: "app", Relation: "authorized-to", Confidence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := graph.AddEdge(Edge{From: "app", To: "db", Relation: "connects-to", Confidence: 1}); err != nil {
		t.Fatal(err)
	}
	paths := graph.Paths("user")
	if len(paths) != 1 || paths[0].Priority != "P0" {
		t.Fatalf("paths=%#v", paths)
	}
}
