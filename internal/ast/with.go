package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&With{})
}

type With struct {
	*InternalNode
	Token      *tokens.Token
	Condition  Node
	Expression Node
}

func (n *With) GetToken() *tokens.Token {
	return n.Token
}

func (n *With) String() string {
	return "<with>"
}

func (n *With) Children() []Node {
	return []Node{n.Condition, n.Expression}
}

func (n *With) Walk(fn WalkFn) {
	n.Condition = fn(n.Condition)
	n.Expression = fn(n.Expression)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
