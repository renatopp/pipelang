package ast

import (
	"encoding/gob"
	"fmt"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&InfixOperator{})
}

type InfixOperator struct {
	*InternalNode
	Token    *tokens.Token
	Operator string
	Left     Node
	Right    Node
}

func (n *InfixOperator) GetToken() *tokens.Token {
	return n.Token
}

func (n *InfixOperator) String() string {
	return fmt.Sprintf("<infix-operator:%s>", n.Operator)
}

func (n *InfixOperator) Children() []Node {
	return []Node{n.Left, n.Right}
}

func (n *InfixOperator) Walk(fn WalkFn) {
	n.Left = fn(n.Left)
	n.Right = fn(n.Right)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
