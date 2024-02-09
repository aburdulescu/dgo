package dgo

type AST struct {
	blocks []any
}

type block struct {
	nodes []any
}

type loopStmt struct {
	stmts []node
}

type altStmt struct {
	stmts     []node
	stmtsElse []node
}

type defStmt struct {
	identifier string
	label      string
}

type msgStmt struct {
	srcIdentifier string
	dstIdentifier string
	label         string
}

type rspStmt struct {
	srcIdentifier string
	dstIdentifier string
	label         string
}
