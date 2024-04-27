package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Unwrap{})
}

type Unwrap struct {
	*InternalNode
	Token  *tokens.Token
	Target Node
}

func (u *Unwrap) GetToken() *tokens.Token {
	return u.Token
}

func (u *Unwrap) String() string {
	return "<unwrap>"
}

func (u *Unwrap) Children() []Node {
	return []Node{u.Target}
}

func (n *Unwrap) Walk(fn WalkFn) {
	n.Target = fn(n.Target)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
