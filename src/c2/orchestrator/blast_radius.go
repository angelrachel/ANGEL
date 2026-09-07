package orchestrator

type BlastRadius struct {
Nodes []string
Edges map[string][]string
}

func NewBlastRadius() *BlastRadius {
return &BlastRadius{
Nodes: []string{},
Edges: make(map[string][]string),
}
}

func (b *BlastRadius) AddNode(node string) {
b.Nodes = append(b.Nodes, node)
}

func (b *BlastRadius) AddEdge(from, to string) {
b.Edges[from] = append(b.Edges[from], to)
}

func (b *BlastRadius) CalculateImpact(node string) int {
visited := make(map[string]bool)
count := 0
var dfs func(current string)
dfs = func(current string) {
visited[current] = true
count++
for _, neighbor := range b.Edges[current] {
if !visited[neighbor] {
dfs(neighbor)
}
}
}
dfs(node)
return count
}

func (b *BlastRadius) GetNodes() []string {
return b.Nodes
}

func (b *BlastRadius) GetEdges() map[string][]string {
return b.Edges
}
