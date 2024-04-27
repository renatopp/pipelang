package ast

import (
	"encoding/gob"
	"fmt"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Number{})
}

type Number struct {
	*InternalNode
	Token *tokens.Token
	Value float64
}

func (n *Number) GetToken() *tokens.Token {
	return n.Token
}

func (n *Number) String() string {
	return fmt.Sprintf("<number:%0.4f>", n.Value)
}

func (n *Number) Children() []Node {
	return []Node{}
}

func (n *Number) Walk(fn WalkFn) {
	//
}
