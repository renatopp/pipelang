package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Break{})
}

type Break struct {
	*InternalNode
	Token *tokens.Token
}

func (n *Break) GetToken() *tokens.Token {
	return n.Token
}

func (n *Break) String() string {
	return "<break>"
}

func (n *Break) Children() []Node {
	return []Node{}
}

func (n *Break) Walk(fn WalkFn) {
	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
