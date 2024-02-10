package dgo

import (
	"fmt"
	"strings"
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

func (n AstNode) Dump() {
	level := 1
	AstWalk(&n, func(v *AstNode) {
		if v.Kind == AstNodeRoot {
			return
		}

		switch v.Kind {
		case AstNodeElse, AstNodeEnd:
			level--
		}

		fmt.Printf("%s%s", strings.Repeat("-", level), v.Kind)
		switch v.Kind {
		case AstNodeDef:
			d := v.Data.(DefStmt)
			fmt.Printf(" %s %q", d.Identifier, d.Label)
		case AstNodeMsg:
			d := v.Data.(MsgStmt)
			fmt.Printf(" %s %s %q", d.Src, d.Dst, d.Label)
		case AstNodeRsp:
			d := v.Data.(RspStmt)
			fmt.Printf(" %s %s %q", d.Src, d.Dst, d.Label)
		case AstNodeAlt:
			d := v.Data.(AltStmt)
			fmt.Printf(" %q", d.Text)
		case AstNodeLoop:
			d := v.Data.(LoopStmt)
			fmt.Printf(" %q", d.Text)
		}
		fmt.Print("\n")

		switch v.Kind {
		case AstNodeRoot, AstNodeLoop, AstNodeAlt, AstNodeElse:
			level++
		}
	})
}

func (n AstNode) String() string {
	return fmt.Sprintf("%s", n.Kind)
}

func (n *AstNode) add(c *AstNode) {
	n.children = append(n.children, c)
}

//go:generate stringer -type=AstNodeKind -trimprefix=AstNode
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
