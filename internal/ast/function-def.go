package ast

import (
	"encoding/gob"
	"fmt"
	"slices"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&FunctionDef{})
}

// Represent a function definition `fn name() {}`
type FunctionDef struct {
	*InternalNode
	Token      *tokens.Token
	Name       string
	Parameters []Node
	Body       Node
	Generator  bool
}

func (n *FunctionDef) GetToken() *tokens.Token {
	return n.Token
}

func (n *FunctionDef) String() string {
	if n.Generator {
		return fmt.Sprintf("<function:%s (generator)>", n.Name)
	}
	return fmt.Sprintf("<function:%s>", n.Name)
}

func (n *FunctionDef) Children() []Node {
	return slices.Concat([]Node{}, n.Parameters, []Node{n.Body})
}

func (n *FunctionDef) Walk(fn WalkFn) {
	for i, child := range n.Parameters {
		n.Parameters[i] = fn(child)
	}
	n.Body = fn(n.Body)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
