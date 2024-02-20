package dgo

import (
	"fmt"
	"strings"
)

type Ast struct {
	nodes []AstNode
}

type AstNode struct {
	Kind AstNodeKind
	Data any
}

func (ast Ast) Walk(fn func(AstNode)) {
	for _, n := range ast.nodes {
		fn(n)
	}
}

func (ast Ast) Dump() {
	level := 1
	for _, node := range ast.nodes {
		switch node.Kind {
		case AstNodeElse, AstNodeEnd:
			level--
		}

		fmt.Printf("%s%s", strings.Repeat("-", level), node.Kind)
		switch node.Kind {
		case AstNodeDef:
			d := node.Data.(DefStmt)
			fmt.Printf(" %s %q", d.Identifier, d.Label)
		case AstNodeMsg:
			d := node.Data.(MsgStmt)
			fmt.Printf(" %s %s %q", d.Src, d.Dst, d.Label)
		case AstNodeRsp:
			d := node.Data.(RspStmt)
			fmt.Printf(" %s %s %q", d.Src, d.Dst, d.Label)
		case AstNodeAlt:
			d := node.Data.(AltStmt)
			fmt.Printf(" %q", d.Text)
		case AstNodeLoop:
			d := node.Data.(LoopStmt)
			fmt.Printf(" %q", d.Text)
		}
		fmt.Print("\n")

		switch node.Kind {
		case AstNodeLoop, AstNodeAlt, AstNodeElse:
			level++
		}
	}
}

func (n AstNode) String() string {
	return fmt.Sprintf("%s", n.Kind)
}

func (ast *Ast) add(node AstNode) {
	ast.nodes = append(ast.nodes, node)
}

func (ast Ast) last() AstNode {
	return ast.nodes[len(ast.nodes)-1]
}

//go:generate stringer -type=AstNodeKind -trimprefix=AstNode
type AstNodeKind uint8

const (
	AstNodeUndefined AstNodeKind = iota
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
