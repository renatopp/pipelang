package ast

import (
	"encoding/gob"
	"fmt"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Identifier{})
}

type Identifier struct {
	*InternalNode
	Token *tokens.Token
	Value string
}

func (n *Identifier) GetToken() *tokens.Token {
	return n.Token
}

func (n *Identifier) String() string {
	return fmt.Sprintf("<identifier:%s>", n.Value)
}

func (n *Identifier) Children() []Node {
	return []Node{}
}

func (n *Identifier) Walk(fn WalkFn) {
	//
}
