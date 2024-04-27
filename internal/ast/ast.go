package ast

import (
	"encoding/gob"

	"github.com/renatopp/langtools/tokens"
)

func init() {
	gob.Register(&tokens.Token{})
}

func PrintTree(node Node) {
	Traverse(node, func(level int, node Node) {
		for i := 0; i < level; i++ {
			print("  ")
		}
		if node == nil {
			println("!ERR: nil node!")
			return
		}
		println(node.String())
	})
}

type TraverseFn func(depth int, node Node)

func Traverse(node Node, f TraverseFn) {
	traverse(0, node, f)
}

func traverse(depth int, node Node, f TraverseFn) {
	if node == nil {
		return
	}

	f(depth, node)

	for _, child := range node.Children() {
		traverse(depth+1, child, f)
	}
}

type WalkFn func(Node) Node

type Node interface {
	GetToken() *tokens.Token
	String() string
	Children() []Node
	SetSourcePath(path string)
	SourcePath() string
	SetParent(parent Node)
	Parent() Node
	Walk(fn WalkFn)
}

type InternalNode struct {
	parent     Node
	sourcePath string
}

func (n *InternalNode) SourcePath() string {
	return n.sourcePath
}

func (n *InternalNode) SetSourcePath(path string) {
	n.sourcePath = path
}

func (n *InternalNode) SetParent(parent Node) {
	n.parent = parent
}

func (n *InternalNode) Parent() Node {
	return n.parent
}
