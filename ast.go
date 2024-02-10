package dgo

import (
	"fmt"
	"strings"
)

type astNode struct {
	kind     astNodeKind
	data     any
	parent   *astNode
	children []*astNode
}

func (n astNode) String() string {
	s := new(strings.Builder)
	fmt.Fprintf(s, "%s:", n.kind)
	for _, c := range n.children {
		fmt.Fprintf(s, " %s", c.kind)
	}
	return s.String()
}

func (n *astNode) add(c *astNode) {
	n.children = append(n.children, c)
}

//go:generate stringer -type astNodeKind
type astNodeKind uint8

const (
	astNodeUndefined astNodeKind = iota
	astNodeRoot
	astNodeDef
	astNodeMsg
	astNodeRsp
	astNodeLoop
	astNodeAlt
	astNodeElse
	astNodeEnd
)

type defStmt struct {
	identifier string
	label      string
}

type msgStmt struct {
	srcIdentifier string
	dstIdentifier string
	label         string
}

type rspStmt struct {
	srcIdentifier string
	dstIdentifier string
	label         string
}
