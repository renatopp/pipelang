package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Dict{})
}

// Represents the dict creation `{a=1, b=2, ...}`
// Elements are a list of [key, value, key, value, ...] pairs
type Dict struct {
	*InternalNode
	Token    *tokens.Token
	Elements []Node
}

func (n *Dict) GetToken() *tokens.Token {
	return n.Token
}

func (n *Dict) String() string {
	return "<dict>"
}

func (n *Dict) Children() []Node {
	return n.Elements
}

func (n *Dict) Walk(fn WalkFn) {
	for i, child := range n.Elements {
		n.Elements[i] = fn(child)
	}

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
