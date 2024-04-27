package ast

import (
	"encoding/gob"
	"fmt"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&PrefixOperator{})
}

type PrefixOperator struct {
	*InternalNode
	Token    *tokens.Token
	Operator string
	Right    Node
}

func (n *PrefixOperator) GetToken() *tokens.Token {
	return n.Token
}

func (n *PrefixOperator) String() string {
	return fmt.Sprintf("<prefix-operator:%s>", n.Operator)
}

func (n *PrefixOperator) Children() []Node {
	return []Node{n.Right}
}

func (n *PrefixOperator) Walk(fn WalkFn) {
	n.Right = fn(n.Right)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
