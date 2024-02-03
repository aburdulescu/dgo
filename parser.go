package dgo

import (
	"fmt"
	"io"

	gotoken "go/token"
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

type parser struct {
	tokens []token
	i      int
	last   token
	nodes  []node
	edges  []edge
	err    error
}

func parse(r io.Reader) error {
	p := parser{tokens: scan(r)}
	return p.parse()
}

func (p *parser) parse() error {
	for !p.token().empty() {
		if p.err != nil {
			return p.err
		}
		p.parseBlock()
	}
	fmt.Println(p.nodes)
	fmt.Println(p.edges)
	return nil
}

func isDef(tok string) bool        { return tok == "def" }
func isMsg(tok string) bool        { return tok == "msg" }
func isRsp(tok string) bool        { return tok == "rsp" }
func isSemicolon(tok string) bool  { return tok == ";" }
func isString(tok string) bool     { return tok[0] == '"' && tok[len(tok)-1] == '"' }
func isIdentifier(tok string) bool { return gotoken.IsIdentifier(tok) }

func (p *parser) parseBlock() {
	switch {
	case p.accept(isDef):
		var n node

		p.expect(isIdentifier)
		n.identifier = p.last.text

		p.expect(isString)
		n.label = p.last.text

		p.expect(isSemicolon)

		p.addNode(n)

	case p.accept(isMsg):
		e := edge{dir: false}

		p.expect(isIdentifier)
		e.src = p.last.text

		p.expect(isIdentifier)
		e.dst = p.last.text

		p.expect(isString)
		e.label = p.last.text

		p.expect(isSemicolon)

		p.addEdge(e)

	case p.accept(isRsp):
		e := edge{dir: true}

		p.expect(isIdentifier)
		e.src = p.last.text

		p.expect(isIdentifier)
		e.dst = p.last.text

		p.expect(isString)
		e.label = p.last.text

		p.expect(isSemicolon)

		p.addEdge(e)

	default:
		p.error("unexpected keyword")
	}
}

func (p *parser) addNode(n node) {
	p.nodes = append(p.nodes, n)
}

func (p *parser) addEdge(e edge) {
	p.edges = append(p.edges, e)
}

func (p *parser) next() {
	p.last = p.token()
	p.i++
}

func (p parser) token() token {
	return p.tokens[p.i]
}

func (p *parser) accept(filter func(string) bool) bool {
	if !filter(p.token().text) {
		return false
	}
	p.next()
	return true
}

func (p *parser) expect(filter func(string) bool) bool {
	if p.err != nil {
		return false
	}
	if p.accept(filter) {
		return true
	}
	p.error("unexpected token")
	return false
}

func (p *parser) error(msg string) {
	tok := p.token()
	p.err = fmt.Errorf("line %d, column %d, text %q: %s", tok.line, tok.column, tok.text, msg)
}
