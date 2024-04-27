package ast

import (
	"encoding/gob"
	"slices"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&If{})
}

type If struct {
	*InternalNode
	Token           *tokens.Token
	Conditions      []Node
	TrueExpression  Node
	FalseExpression Node
}

func (n *If) GetToken() *tokens.Token {
	return n.Token
}

func (n *If) String() string {
	return "<if>"
}

func (n *If) Children() []Node {
	return slices.Concat(n.Conditions, []Node{n.TrueExpression, n.FalseExpression})
}

func (n *If) Walk(fn WalkFn) {
	for i := range n.Conditions {
		n.Conditions[i] = fn(n.Conditions[i])
	}
	n.TrueExpression = fn(n.TrueExpression)
	n.FalseExpression = fn(n.FalseExpression)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
