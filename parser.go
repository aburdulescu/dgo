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
	d    Diagram
}

func Parse(r io.Reader) (*Diagram, error) {
	p := parser{t: scan(r)}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return &p.d, nil
}

func (p *parser) parse() error {
	for !p.token().empty() {
		if p.err != nil {
			return p.err
		}
		p.parseBlock()
	}
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

		p.d.addNode(n)

	case p.accept(isMsg):
		e := edge{}

		p.expect(isIdentifier)
		e.src = p.d.find(p.last.text)
		if e.src == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}

		p.expect(isIdentifier)
		e.dst = p.d.find(p.last.text)
		if e.dst == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}

		p.expect(isString)
		e.label = stripString(p.last.text)

		p.expect(isSemicolon)

		p.d.addMsg(e)

	case p.accept(isRsp):
		e := edge{}

		p.expect(isIdentifier)
		e.src = p.d.find(p.last.text)
		if e.src == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}

		p.expect(isIdentifier)
		e.dst = p.d.find(p.last.text)
		if e.dst == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}

		p.expect(isString)
		e.label = stripString(p.last.text)

		p.expect(isSemicolon)

		p.d.addRsp(e)

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
