package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Wrap{})
}

type Wrap struct {
	*InternalNode
	Token  *tokens.Token
	Target Node
}

func (u *Wrap) GetToken() *tokens.Token {
	return u.Token
}

func (u *Wrap) String() string {
	return "<wrap>"
}

func (u *Wrap) Children() []Node {
	return []Node{u.Target}
}

func (n *Wrap) Walk(fn WalkFn) {
	n.Target = fn(n.Target)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
