package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Yield{})
}

type Yield struct {
	*InternalNode
	Token      *tokens.Token
	Expression Node
	Break      bool
}

func (n *Yield) GetToken() *tokens.Token {
	return n.Token
}

func (n *Yield) String() string {
	if n.Break {
		return "<yield:break>"
	}

	return "<yield>"
}

func (n *Yield) Children() []Node {
	if n.Expression == nil {
		return []Node{}
	}
	return []Node{n.Expression}
}

func (n *Yield) Walk(fn WalkFn) {
	n.Expression = fn(n.Expression)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
