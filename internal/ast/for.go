package ast

import (
	"encoding/gob"
	"slices"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&For{})
}

type For struct {
	*InternalNode
	Token        *tokens.Token
	Conditions   []Node
	InExpression Node
	Expression   Node
}

func (n *For) GetToken() *tokens.Token {
	return n.Token
}

func (n *For) String() string {
	if n.InExpression != nil {
		return "<for-in>"
	}
	return "<for>"
}

func (n *For) Children() []Node {
	if n.InExpression != nil {
		return slices.Concat(n.Conditions, []Node{n.InExpression, n.Expression})
	}
	return slices.Concat(n.Conditions, []Node{n.Expression})
}

func (n *For) Walk(fn WalkFn) {
	for i := range n.Conditions {
		n.Conditions[i] = fn(n.Conditions[i])
	}
	if n.InExpression != nil {
		n.InExpression = fn(n.InExpression)
	}
	n.Expression = fn(n.Expression)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
