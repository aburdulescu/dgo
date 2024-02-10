package dgo

import (
	"fmt"
	"io"

	gotoken "go/token"
)

type parser struct {
	t            []token
	i            int
	last         token
	err          error
	participants []string
	ast          *AstNode
}

func (p *parser) addParticipant(id string) {
	p.participants = append(p.participants, id)
}

func (p parser) findParticipant(id string) int {
	for i, v := range p.participants {
		if v == id {
			return i
		}
	}
	return -1
}

func Parse(r io.Reader) (*AstNode, error) {
	root := &AstNode{Kind: AstNodeRoot}
	p := parser{
		t:   scan(r),
		ast: root,
	}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return root, nil
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
		var d DefStmt

		p.expect(isIdentifier)
		d.Identifier = p.last.text

		p.expect(isEqual)

		p.expect(isString)
		d.Label = stripString(p.last.text)

		p.addParticipant(d.Identifier)

		p.ast.add(&AstNode{
			Kind:   AstNodeDef,
			Data:   d,
			parent: p.ast,
		})

	case p.accept(isMsg):
		var d MsgStmt

		p.expect(isIdentifier)
		if i := p.findParticipant(p.last.text); i == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}
		d.Src = p.last.text

		p.expect(isDash)
		p.expect(isLess)

		p.expect(isIdentifier)
		if i := p.findParticipant(p.last.text); i == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}
		d.Dst = p.last.text

		p.expect(isEqual)

		p.expect(isString)
		d.Label = stripString(p.last.text)

		p.ast.add(&AstNode{
			Kind:   AstNodeMsg,
			Data:   d,
			parent: p.ast,
		})

	case p.accept(isRsp):
		var d RspStmt

		p.expect(isIdentifier)
		if i := p.findParticipant(p.last.text); i == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}
		d.Src = p.last.text

		p.expect(isDash)
		p.expect(isLess)

		p.expect(isIdentifier)
		if i := p.findParticipant(p.last.text); i == -1 {
			p.error(fmt.Sprintf("unknown identifier '%s'", p.last.text))
			return
		}
		d.Dst = p.last.text

		p.expect(isEqual)

		p.expect(isString)
		d.Label = stripString(p.last.text)

		p.ast.add(&AstNode{
			Kind:   AstNodeRsp,
			Data:   d,
			parent: p.ast,
		})

	case p.accept(isAlt):
		var d AltStmt
		p.expect(isString)
		d.Text = stripString(p.last.text)

		n := &AstNode{
			Kind:   AstNodeAlt,
			Data:   d,
			parent: p.ast,
		}
		n.parent.add(n)
		p.ast = n

	case p.accept(isElse):
		if p.ast.Kind != AstNodeAlt {
			// error
		}

		n := &AstNode{
			Kind:   AstNodeElse,
			parent: p.ast.parent,
		}
		n.parent.add(n)
		p.ast = n

	case p.accept(isEnd):
		if p.ast.Kind != AstNodeElse &&
			p.ast.Kind != AstNodeLoop {
			// error
		}

		n := &AstNode{
			Kind:   AstNodeEnd,
			parent: p.ast.parent,
		}
		n.parent.add(n)
		p.ast = n.parent

	case p.accept(isLoop):
		var d LoopStmt

		p.expect(isString)
		d.Text = stripString(p.last.text)

		n := &AstNode{
			Kind:   AstNodeLoop,
			Data:   d,
			parent: p.ast,
		}
		n.parent.add(n)
		p.ast = n

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
