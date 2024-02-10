package dgo

import (
	"fmt"
)

type AstNode struct {
	Kind     AstNodeKind
	Data     any
	parent   *AstNode
	children []*AstNode
}

func AstWalk(n *AstNode, fn func(*AstNode)) {
	fn(n)
	for _, c := range n.children {
		AstWalk(c, fn)
	}
}

func (n AstNode) String() string {
	return fmt.Sprintf("%s", n.Kind)
}

func (n *AstNode) add(c *AstNode) {
	n.children = append(n.children, c)
}

//go:generate stringer -type AstNodeKind
type AstNodeKind uint8

const (
	AstNodeUndefined AstNodeKind = iota
	AstNodeRoot
	AstNodeDef
	AstNodeMsg
	AstNodeRsp
	AstNodeLoop
	AstNodeAlt
	AstNodeElse
	AstNodeEnd
)

type DefStmt struct {
	Identifier string
	Label      string
}

type MsgStmt struct {
	Src   string
	Dst   string
	Label string
}

type RspStmt struct {
	Src   string
	Dst   string
	Label string
}

type AltStmt struct {
	Text string
}

type LoopStmt struct {
	Text string
}
