package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Block{})
}

type Block struct {
	*InternalNode
	Token       *tokens.Token
	Expressions []Node
	Scoped      bool
}

func (n *Block) GetToken() *tokens.Token {
	return n.Token
}

func (n *Block) String() string {
	if !n.Scoped {
		return "<block:non-scoped>"
	}
	return "<block>"
}

func (n *Block) Children() []Node {
	return n.Expressions
}

func (n *Block) Walk(fn WalkFn) {
	for i, child := range n.Expressions {
		n.Expressions[i] = fn(child)
	}

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
