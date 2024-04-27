package ast

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&String{})
}

type String struct {
	*InternalNode
	Token *tokens.Token
	Value string
}

func (n *String) GetToken() *tokens.Token {
	return n.Token
}

func (n *String) String() string {
	value := strings.ReplaceAll(n.Value, "\n", "\\n")
	return fmt.Sprintf("<string:%s>", value)
}

func (n *String) Children() []Node {
	return []Node{}
}

func (n *String) Walk(fn WalkFn) {
	//
}
