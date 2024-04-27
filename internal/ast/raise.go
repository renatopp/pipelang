package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Raise{})
}

type Raise struct {
	*InternalNode
	Token      *tokens.Token
	Expression Node
}

func (n *Raise) GetToken() *tokens.Token {
	return n.Token
}

func (n *Raise) String() string {
	return "<raise>"
}

func (n *Raise) Children() []Node {
	return []Node{n.Expression}
}

func (n *Raise) Walk(fn WalkFn) {
	n.Expression = fn(n.Expression)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
