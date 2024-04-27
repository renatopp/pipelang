package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Access{})
}

type Access struct {
	*InternalNode
	Token *tokens.Token
	Left  Node
	Right Node
}

func (n *Access) GetToken() *tokens.Token {
	return n.Token
}

func (n *Access) String() string {
	return "<access>"
}

func (n *Access) Children() []Node {
	return []Node{n.Left, n.Right}
}

func (n *Access) Walk(fn WalkFn) {
	n.Left = fn(n.Left)
	n.Right = fn(n.Right)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
