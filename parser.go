package dgo

import (
	"fmt"
	"io"

	gotoken "go/token"
)

type parser struct {
	t          []token
	i          int
	last       token
	err        error
	d          Diagram
	astRoot    *astNode
	astCurrent *astNode
}

func Parse(r io.Reader) (*Diagram, *astNode, error) {
	root := &astNode{kind: astNodeRoot}
	p := parser{
		t:          scan(r),
		astRoot:    root,
		astCurrent: root,
	}
	if err := p.parse(); err != nil {
		return nil, nil, err
	}
	return &p.d, p.astRoot, nil
}

func (p *parser) parse() error {
	for !p.token().empty() {
		if p.err != nil {
			return p.err
		}
		p.parseStmt()
	}
	return nil
}

func isDef(tok string) bool          { return tok == "def" }
func isMsg(tok string) bool          { return tok == "msg" }
func isRsp(tok string) bool          { return tok == "rsp" }
func isRawString(tok string) bool    { return tok[0] == '`' && tok[len(tok)-1] == '`' }
func isSimpleString(tok string) bool { return tok[0] == '"' && tok[len(tok)-1] == '"' }
func isString(tok string) bool       { return isRawString(tok) || isSimpleString(tok) }
func isIdentifier(tok string) bool   { return gotoken.IsIdentifier(tok) }
func isAlt(tok string) bool          { return tok == "alt" }
func isElse(tok string) bool         { return tok == "else" }
func isEnd(tok string) bool          { return tok == "end" }
func isLoop(tok string) bool         { return tok == "loop" }
func isEqual(tok string) bool        { return tok == "=" }
func isDash(tok string) bool         { return tok == "-" }
func isLess(tok string) bool         { return tok == ">" }

func stripString(s string) string {
	return s[1 : len(s)-1]
}

func (p *parser) parseStmt() {
	switch {
	case p.accept(isDef):
		var n node

		p.expect(isIdentifier)
		n.identifier = p.last.text

		p.expect(isEqual)

		p.expect(isString)
		n.label = stripString(p.last.text)

		p.d.addNode(n)

		p.astCurrent.add(&astNode{
			kind: astNodeDef,
			data: defStmt{
				identifier: n.identifier,
				label:      n.label,
			},
		})

	case p.accept(isMsg):
		e := edge{}

		p.expect(isIdentifier)
		e.src = p.d.find(p.last.text)
		if e.src == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}

		p.expect(isDash)
		p.expect(isLess)

		p.expect(isIdentifier)
		e.dst = p.d.find(p.last.text)
		if e.dst == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}

		p.expect(isEqual)

		p.expect(isString)
		e.label = stripString(p.last.text)

		p.d.addMsg(e)

		p.astCurrent.add(&astNode{
			kind: astNodeMsg,
			data: msgStmt{
				srcIdentifier: p.d.nodes[e.src].identifier,
				dstIdentifier: p.d.nodes[e.dst].identifier,
				label:         e.label,
			},
		})

	case p.accept(isRsp):
		e := edge{}

		p.expect(isIdentifier)
		e.src = p.d.find(p.last.text)
		if e.src == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}

		p.expect(isDash)
		p.expect(isLess)

		p.expect(isIdentifier)
		e.dst = p.d.find(p.last.text)
		if e.dst == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}

		p.expect(isEqual)

		p.expect(isString)
		e.label = stripString(p.last.text)

		p.d.addRsp(e)

		p.astCurrent.add(&astNode{
			kind: astNodeRsp,
			data: rspStmt{
				srcIdentifier: p.d.nodes[e.src].identifier,
				dstIdentifier: p.d.nodes[e.dst].identifier,
				label:         e.label,
			},
		})

	case p.accept(isAlt):
		p.expect(isString)

		n := &astNode{
			kind:   astNodeAlt,
			parent: p.astCurrent,
		}
		n.parent.add(n)
		p.astCurrent = n

	case p.accept(isElse):
		if p.astCurrent.kind != astNodeAlt {
			// error
		}

		n := &astNode{
			kind:   astNodeElse,
			parent: p.astCurrent.parent,
		}
		n.parent.add(n)
		p.astCurrent = n

	case p.accept(isEnd):
		if p.astCurrent.kind != astNodeElse &&
			p.astCurrent.kind != astNodeLoop {
			// error
		}

		n := &astNode{
			kind:   astNodeEnd,
			parent: p.astCurrent.parent,
		}
		n.parent.add(n)
		p.astCurrent = n.parent

	case p.accept(isLoop):
		p.expect(isString)

		n := &astNode{
			kind:   astNodeLoop,
			parent: p.astCurrent,
		}
		n.parent.add(n)
		p.astCurrent = n

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
