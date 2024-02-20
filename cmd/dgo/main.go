package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aburdulescu/dgo"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func mainErr() error {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: dgo <command>

Commands:
    fmt  Format the code read from stdin
    ast  Dump the AST for the code read from stdin

`)
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return nil
	}

	ast, err := dgo.Parse(os.Stdin)
	if err != nil {
		return err
	}

	switch cmd := flag.Arg(0); cmd {
	case "fmt":
		format(ast)
	case "ast":
		ast.Dump()
	default:
		return fmt.Errorf("unknown command '%s'", cmd)
	}

	return nil
}

func format(ast *dgo.Ast) {
	// TODO: move to fmt pkg/cmd
	// TODO: preserve newlines and comment
	// TODO: preserve raw strings
	indent := 0
	ast.Walk(func(n dgo.AstNode) {
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
