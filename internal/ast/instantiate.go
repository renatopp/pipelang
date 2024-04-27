package ast

import (
	"encoding/gob"
	"slices"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Instantiate{})
}

// Represents an instantiation of a type. `CustomType { x=2 }`
type Instantiate struct {
	*InternalNode
	Token    *tokens.Token
	Target   Node
	Elements []Node
}

func (n *Instantiate) GetToken() *tokens.Token {
	return n.Token
}

func (n *Instantiate) String() string {
	return "<instantiate>"
}

func (n *Instantiate) Children() []Node {
	return slices.Concat([]Node{n.Target}, n.Elements)
}

func (n *Instantiate) Walk(fn WalkFn) {
	n.Target = fn(n.Target)
	for i, child := range n.Elements {
		n.Elements[i] = fn(child)
	}

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
