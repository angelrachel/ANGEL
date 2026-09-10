package orchestrator

type AttackGraph struct {
	Nodes []string
	Edges map[string][]string
}

func NewAttackGraph() *AttackGraph {
	return &AttackGraph{
		Nodes: []string{},
		Edges: make(map[string][]string),
	}
}

func (a *AttackGraph) AddNode(node string) {
	a.Nodes = append(a.Nodes, node)
}

func (a *AttackGraph) AddEdge(from, to string) {
	a.Edges[from] = append(a.Edges[from], to)
}

func (a *AttackGraph) GetPaths(from, to string) [][]string {
	var paths [][]string
	var currentPath []string
	visited := make(map[string]bool)
	var dfs func(node string)
	dfs = func(node string) {
		visited[node] = true
		currentPath = append(currentPath, node)
		if node == to {
			pathCopy := make([]string, len(currentPath))
			copy(pathCopy, currentPath)
			paths = append(paths, pathCopy)
		} else {
			for _, neighbor := range a.Edges[node] {
				if !visited[neighbor] {
					dfs(neighbor)
				}
			}
		}
		currentPath = currentPath[:len(currentPath)-1]
		visited[node] = false
	}
	dfs(from)
	return paths
}

func (a *AttackGraph) GetNodes() []string {
	return a.Nodes
}

func (a *AttackGraph) GetEdges() map[string][]string {
	return a.Edges
}
