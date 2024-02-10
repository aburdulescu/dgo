package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aburdulescu/dgo"
)

func main() {
	a, err := dgo.Parse(os.Stdin)
	if err != nil {
		panic(err)
	}

	// TODO: move to fmt pkg/cmd
	// TODO: preserve newlines and comment
	// TODO: preserve raw strings

	indent := 0

	dgo.AstWalk(a, func(n *dgo.AstNode) {
		if n.Kind == dgo.AstNodeRoot {
			return
		}

		switch n.Kind {
		case dgo.AstNodeElse, dgo.AstNodeEnd:
			indent -= 2
		}

		switch n.Kind {
		case dgo.AstNodeDef:
			data := n.Data.(dgo.DefStmt)
			fmt.Printf("%sdef %s = %q\n", strings.Repeat(" ", indent), data.Identifier, data.Label)

		case dgo.AstNodeMsg:
			data := n.Data.(dgo.MsgStmt)
			fmt.Printf("%smsg %s -> %s = %q\n", strings.Repeat(" ", indent), data.Src, data.Dst, data.Label)

		case dgo.AstNodeRsp:
			data := n.Data.(dgo.RspStmt)
			fmt.Printf("%srsp %s -> %s = %q\n", strings.Repeat(" ", indent), data.Src, data.Dst, data.Label)

		case dgo.AstNodeAlt:
			data := n.Data.(dgo.AltStmt)
			fmt.Printf("%salt %q\n", strings.Repeat(" ", indent), data.Text)

		case dgo.AstNodeElse:
			fmt.Printf("%selse\n", strings.Repeat(" ", indent))

		case dgo.AstNodeLoop:
			data := n.Data.(dgo.LoopStmt)
			fmt.Printf("%sloop %q\n", strings.Repeat(" ", indent), data.Text)

		case dgo.AstNodeEnd:
			fmt.Printf("%send\n", strings.Repeat(" ", indent))

		default:
			fmt.Printf("%s%s\n", strings.Repeat(" ", indent), n)
		}

		switch n.Kind {
		case dgo.AstNodeLoop, dgo.AstNodeAlt, dgo.AstNodeElse:
			indent += 2
		}
	})
}
