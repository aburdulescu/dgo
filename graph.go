package dgo

import (
	"fmt"
	"strings"
)

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

func (g graph) Write() string {
	s := new(strings.Builder)
	fmt.Fprintf(s, "digraph {\n")
	for _, node := range g.nodes {
		fmt.Fprintf(s, "  %s_s [shape=box, label=\"%s\"];\n", node.identifier, node.label)
		fmt.Fprintf(s, "  %s_e [shape=box,label=\"%s\"];\n", node.identifier, node.label)
		fmt.Fprintf(s, "  %s_s -> %s_e [arrowhead=none];\n", node.identifier, node.identifier)
	}
	fmt.Fprintf(s, "}")
	return s.String()
}
