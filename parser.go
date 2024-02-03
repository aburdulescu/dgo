package dgo

import (
	"fmt"
	"io"

	gotoken "go/token"
)

type parser struct {
	t    []token
	i    int
	last token
	err  error
	g    graph
}

func parse(r io.Reader) (*graph, error) {
	p := parser{t: scan(r)}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return &p.g, nil
}

func (p *parser) parse() error {
	for !p.token().empty() {
		if p.err != nil {
			return p.err
		}
		p.parseBlock()
	}
	fmt.Println(p.g.nodes)
	fmt.Println(p.g.edges)
	return nil
}

func isDef(tok string) bool          { return tok == "def" }
func isMsg(tok string) bool          { return tok == "msg" }
func isRsp(tok string) bool          { return tok == "rsp" }
func isSemicolon(tok string) bool    { return tok == ";" }
func isRawString(tok string) bool    { return tok[0] == '`' && tok[len(tok)-1] == '`' }
func isSimpleString(tok string) bool { return tok[0] == '"' && tok[len(tok)-1] == '"' }
func isString(tok string) bool       { return isRawString(tok) || isSimpleString(tok) }
func isIdentifier(tok string) bool   { return gotoken.IsIdentifier(tok) }

func stripString(s string) string {
	return s[1 : len(s)-1]
}

func (p *parser) parseBlock() {
	switch {
	case p.accept(isDef):
		var n node

		p.expect(isIdentifier)
		n.identifier = p.last.text

		p.expect(isString)
		n.label = stripString(p.last.text)

		p.expect(isSemicolon)

		p.g.addNode(n)

	case p.accept(isMsg):
		e := edge{dir: false}

		p.expect(isIdentifier)
		e.src = p.last.text

		p.expect(isIdentifier)
		e.dst = p.last.text

		p.expect(isString)
		e.label = stripString(p.last.text)

		p.expect(isSemicolon)

		p.g.addEdge(e)

	case p.accept(isRsp):
		e := edge{dir: true}

		p.expect(isIdentifier)
		e.src = p.last.text

		p.expect(isIdentifier)
		e.dst = p.last.text

		p.expect(isString)
		e.label = stripString(p.last.text)

		p.expect(isSemicolon)

		p.g.addEdge(e)

	default:
		p.error("unexpected keyword")
	}
}

func (p *parser) next() {
	p.last = p.token()
	p.i++
}

func (p parser) token() token {
	return p.t[p.i]
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
