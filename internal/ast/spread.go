package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Spread{})
}

type Spread struct {
	*InternalNode
	Token  *tokens.Token
	Target Node
	In     bool
}

func (u *Spread) GetToken() *tokens.Token {
	return u.Token
}

func (u *Spread) String() string {
	if u.In {
		return "<spread:in>"
	}
	return "<spread:out>"
}

func (u *Spread) Children() []Node {
	return []Node{u.Target}
}

func (n *Spread) Walk(fn WalkFn) {
	n.Target = fn(n.Target)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
