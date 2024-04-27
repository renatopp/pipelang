package ast

import (
	"encoding/gob"
	"fmt"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&DataDef{})
}

// Represent a function definition `fn name() {}`
type DataDef struct {
	*InternalNode
	Token      *tokens.Token
	Name       string
	Extensions []Node
	Attributes map[string]Node
	Methods    map[string]Node
}

func (n *DataDef) GetToken() *tokens.Token {
	return n.Token
}

func (n *DataDef) String() string {
	return fmt.Sprintf("<data:%s>", n.Name)
}

func (n *DataDef) Children() []Node {
	children := make([]Node, 0, len(n.Extensions)+len(n.Attributes)+len(n.Methods))
	children = append(children, n.Extensions...)

	for _, child := range n.Attributes {
		children = append(children, child)
	}

	for _, child := range n.Methods {
		children = append(children, child)
	}

	return children
}

func (n *DataDef) Walk(fn WalkFn) {
	for i, child := range n.Extensions {
		n.Extensions[i] = fn(child)
	}

	for i, child := range n.Attributes {
		n.Attributes[i] = fn(child)
	}

	for i, child := range n.Methods {
		n.Methods[i] = fn(child)
	}

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
