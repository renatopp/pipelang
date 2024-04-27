package ast

import (
	"encoding/gob"
	"fmt"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Assignment{})
}

type Assignment struct {
	*InternalNode
	Token    *tokens.Token
	Operator string
	Left     Node
	Right    Node
}

func (n *Assignment) GetToken() *tokens.Token {
	return n.Token
}

func (n *Assignment) String() string {
	return fmt.Sprintf("<assignment:%s>", n.Operator)
}

func (n *Assignment) Children() []Node {
	return []Node{n.Left, n.Right}
}

func (n *Assignment) Walk(fn WalkFn) {
	n.Left = fn(n.Left)
	n.Right = fn(n.Right)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
