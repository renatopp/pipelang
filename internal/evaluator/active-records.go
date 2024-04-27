package evaluator

import o "github.com/renatopp/pipelang/internal/object"

type BlockRecord struct {
	Scope     *o.Scope
	Statement int
}

type IfRecord struct {
	Scope     *o.Scope
	Condition bool
}

type ForRecord struct {
	Scope *o.Scope
}

type WithRecord struct {
	Scope *o.Scope
}

type MatchRecord struct {
	Scope     *o.Scope
	Case      int
	CaseScope *o.Scope
}
