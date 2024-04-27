package ast

import (
	"encoding/gob"
	"slices"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Match{})
}

type Match struct {
	*InternalNode
	Token      *tokens.Token
	Expression Node
	Cases      []Node // [condition, expression, condition, expression, ...]
}

func (n *Match) GetToken() *tokens.Token {
	return n.Token
}

func (n *Match) String() string {
	return "<match>"
}

func (n *Match) Children() []Node {
	return slices.Concat([]Node{n.Expression}, n.Cases)
}

func (n *Match) Walk(fn WalkFn) {
	n.Expression = fn(n.Expression)
	for i := range n.Cases {
		n.Cases[i] = fn(n.Cases[i])
	}

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
