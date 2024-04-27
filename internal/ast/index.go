package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&Index{})
}

// Represents access like `a[0]` and `a[a, b]`
type Index struct {
	*InternalNode
	Token  *tokens.Token
	Target Node
	Index  Node
}

func (n *Index) GetToken() *tokens.Token {
	return n.Token
}

func (n *Index) String() string {
	return "<index>"
}

func (n *Index) Children() []Node {
	return []Node{n.Target, n.Index}
}

func (n *Index) Walk(fn WalkFn) {
	n.Target = fn(n.Target)
	n.Index = fn(n.Index)

	for _, child := range n.Children() {
		child.Walk(fn)
	}
}
