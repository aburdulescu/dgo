package dgo

type AST struct {
	stmts []any
}

type loopStmt struct {
	stmts []any
}

type altStmt struct {
	stmts     []any
	stmtsElse []any
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
