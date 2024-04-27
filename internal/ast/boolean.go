package ast

import (
	"encoding/gob"
	"fmt"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Boolean{})
}

type Boolean struct {
	*InternalNode
	Token *tokens.Token
	Value bool
}

func (n *Boolean) GetToken() *tokens.Token {
	return n.Token
}

func (n *Boolean) String() string {
	return fmt.Sprintf("<boolean:%v>", n.Value)
}

func (n *Boolean) Children() []Node {
	return []Node{}
}

func (n *Boolean) Walk(fn WalkFn) {
	//
}
