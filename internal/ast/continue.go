package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Continue{})
}

type Continue struct {
	*InternalNode
	Token *tokens.Token
}

func (n *Continue) GetToken() *tokens.Token {
	return n.Token
}

func (n *Continue) String() string {
	return "<continue>"
}

func (n *Continue) Children() []Node {
	return []Node{}
}

func (n *Continue) Walk(fn WalkFn) {
	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
