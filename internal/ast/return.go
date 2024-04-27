package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Return{})
}

type Return struct {
	*InternalNode
	Token      *tokens.Token
	Expression Node
}

func (n *Return) GetToken() *tokens.Token {
	return n.Token
}

func (n *Return) String() string {
	return "<return>"
}

func (n *Return) Children() []Node {
	return []Node{n.Expression}
}

func (n *Return) Walk(fn WalkFn) {
	n.Expression = fn(n.Expression)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
