package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Tuple{})
}

type Tuple struct {
	*InternalNode
	Token    *tokens.Token
	Elements []Node
}

func (n *Tuple) GetToken() *tokens.Token {
	return n.Token
}

func (n *Tuple) String() string {
	return "<tuple>"
}

func (n *Tuple) Children() []Node {
	return n.Elements
}

func (n *Tuple) Walk(fn WalkFn) {
	for i, child := range n.Elements {
		n.Elements[i] = fn(child)
	}

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
