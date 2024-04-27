package ast

import (
	"encoding/gob"
	"slices"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Call{})
}

// Represent a function or type call, i.e., any `identifier(<expression>, ...)`
type Call struct {
	*InternalNode
	Token     *tokens.Token
	Target    Node
	Arguments []Node
}

func (n *Call) GetToken() *tokens.Token {
	return n.Token
}

func (n *Call) String() string {
	return "<call>"
}

func (n *Call) Children() []Node {
	return slices.Concat([]Node{n.Target}, n.Arguments)
}

func (n *Call) Walk(fn WalkFn) {
	n.Target = fn(n.Target)
	for i, child := range n.Arguments {
		n.Arguments[i] = fn(child)
	}

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
