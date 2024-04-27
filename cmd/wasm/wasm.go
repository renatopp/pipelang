//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/renatopp/pipelang/internal/ast"
	"github.com/renatopp/pipelang/internal/object"
	"github.com/renatopp/pipelang/internal/runtime"
)

func main() {
	rt := runtime.New()

	js.Global().Set("pipe_eval", js.FuncOf(func(this js.Value, args []js.Value) any {
		code := args[0].String()

		node, err := rt.LoadAst([]byte(code))
		if err != nil {
			return js.ValueOf(err.Error())
		}

		block := node.(*ast.Block)
		var res object.Object
		for _, exp := range block.Expressions {
			res, err = rt.RunAst(exp)
			if err != nil {
				return js.ValueOf(err.Error())
			}
		}

		return js.ValueOf(res.AsString())
	}))

	js.Global().Set("pipe_run", js.FuncOf(func(this js.Value, args []js.Value) any {
		rt := runtime.New()
		code := args[0].String()

		node, err := rt.LoadAst([]byte(code))
		if err != nil {
			return js.ValueOf(err.Error())
		}

		block := node.(*ast.Block)
		var res object.Object
		for _, exp := range block.Expressions {
			res, err = rt.RunAst(exp)
			if err != nil {
				return js.ValueOf(err.Error())
			}
		}

		return js.ValueOf(res.AsString())
	}))

	select {} // keep running
}
