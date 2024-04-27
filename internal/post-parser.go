package internal

import "github.com/renatopp/pipelang/internal/ast"

func OptimizeAst(root ast.Node) ast.Node {
	root = walker(root)
	root.Walk(walker)
	return root
}

func walker(node ast.Node) ast.Node {
	switch node := node.(type) {

	// Force tuple on left side of assignment ===> a = 2 -> (a) = 2
	// Force expansion of assignment operators ===> a += 2 -> (a) = a + 2
	case *ast.Assignment:
		if node.Operator == "=" || node.Operator == ":=" {
			node.Left = checkLeftSideTuple(node.Left)
			return node
		}

		op := node.Operator[:len(node.Operator)-1]

		return &ast.Assignment{
			Token:    node.Token,
			Operator: "=",
			Left:     checkLeftSideTuple(node.Left),
			Right: &ast.InfixOperator{
				Token:    node.Token,
				Left:     node.Left,
				Operator: op,
				Right:    node.Right,
			},
		}
	}

	return node
}

// Force the left side to be a tuple
func checkLeftSideTuple(left ast.Node) ast.Node {
	if _, ok := left.(*ast.Tuple); !ok {
		return &ast.Tuple{
			Token:    left.GetToken(),
			Elements: []ast.Node{left},
		}
	}

	return left
}
