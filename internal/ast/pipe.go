package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Pipe{})
}

type Pipe struct {
	*InternalNode
	Token *tokens.Token
	Left  Node
	Right Node
}

func (n *Pipe) GetToken() *tokens.Token {
	return n.Token
}

func (n *Pipe) String() string {
	return "<pipe>"
}

func (n *Pipe) Children() []Node {
	return []Node{n.Left, n.Right}
}

func (n *Pipe) Walk(fn WalkFn) {
	n.Left = fn(n.Left)
	n.Right = fn(n.Right)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
