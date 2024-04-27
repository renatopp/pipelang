package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&List{})
}

// Represents the list creation `[1, 2, ...]`
type List struct {
	*InternalNode
	Token    *tokens.Token
	Elements []Node
}

func (n *List) GetToken() *tokens.Token {
	return n.Token
}

func (n *List) String() string {
	return "<list>"
}

func (n *List) Children() []Node {
	return n.Elements
}

func (n *List) Walk(fn WalkFn) {
	for i, child := range n.Elements {
		n.Elements[i] = fn(child)
	}

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
