package internal

import "github.com/renatopp/pipelang/internal/ast"

type Error struct {
	Message string
	Node    ast.Node
	Stack   []ast.Node
}

func NewError(message string, stack []ast.Node) *Error {
	var node ast.Node
	if len(stack) > 0 {
		node = stack[len(stack)-1]
	}

	return &Error{
		Message: message,
		Node:    node,
		Stack:   stack,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) At() (line, column int) {
	if e.Node == nil {
		return 0, 0
	}

	minLine, minColumn := e.Node.GetToken().From()
	stack := []ast.Node{e.Node}
	for len(stack) > 0 {
		node := stack[0]
		stack = append(stack[1:], node.Children()...)

		l, c := node.GetToken().From()
		minLine = min(minLine, l)
		minColumn = min(minColumn, c)
	}

	return minLine, minColumn
}

func (e *Error) Range() (fromLine, fromCol, toLine, toCol int) {
	if e.Node == nil {
		return 0, 0, 0, 0
	}

	fromLine, fromCol = e.Node.GetToken().From()
	toLine, toCol = e.Node.GetToken().To()

	stack := []ast.Node{e.Node}
	for len(stack) > 0 {
		node := stack[0]
		stack = append(stack[1:], node.Children()...)

		l, c := node.GetToken().From()
		fromLine = min(fromLine, l)
		fromCol = min(fromCol, c)

		l, c = node.GetToken().To()
		toLine = max(toLine, l)
		toCol = max(toCol, c)
	}

	return fromLine, fromCol, toLine, toCol
}
