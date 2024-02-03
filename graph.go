package dgo

type node struct {
	identifier string
	label      string
}

type edge struct {
	src, dst string
	dir      bool // msg=false, rsp=true
	label    string
}

type graph struct {
	nodes []node
	edges []edge
}

func (g *graph) addNode(n node) {
	g.nodes = append(g.nodes, n)
}

func (g *graph) addEdge(e edge) {
	g.edges = append(g.edges, e)
}
