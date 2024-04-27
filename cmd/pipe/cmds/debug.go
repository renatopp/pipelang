package cmds

import (
	"fmt"
	"os"

	"github.com/renatopp/pipelang/internal"
	"github.com/renatopp/pipelang/internal/ast"
	"github.com/renatopp/pipelang/internal/builtins"
	"github.com/renatopp/pipelang/internal/errfmt"
	"github.com/renatopp/pipelang/internal/evaluator"
	"github.com/renatopp/pipelang/internal/object"
)

func Debug() {
	source := []byte(`
	(a, ...c, d) := [1,2,3,4]...;
	`)

	if len(os.Args) > 2 {
		source = []byte(os.Args[2])
	}

	debugTokenizer(source)
}

func debugTokenizer(source []byte) {
	println("tokenization step:")

	preLexer := internal.NewPreLexer(source)
	tokens := preLexer.All()
	if preLexer.HasErrors() {
		err := errfmt.FormatLexerErrors(preLexer.Errors(), source, "<stdin>")
		fmt.Println(err)
		os.Exit(1)
	}

	for _, token := range tokens {
		fmt.Println(token.DebugString())
	}

	println("\npost-lexing step:")
	postLexer := internal.NewPostLexer(tokens)
	if err := postLexer.Optimize(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tokens = postLexer.All()

	for _, token := range tokens {
		fmt.Println(token.DebugString())
	}

	println("\nparser step:")
	parser := internal.NewPipeParser(postLexer)
	// parser.Log = log.New(os.Stdout, "", 0)
	root := parser.Parse()
	if parser.HasErrors() {
		err := errfmt.FormatParserErrors(parser.Errors(), source, "<stdin>")
		fmt.Println(err)
		os.Exit(1)
	}
	ast.PrintTree(root)

	println("\noptimization step:")
	root = internal.OptimizeAst(root)
	ast.PrintTree(root)

	println("\nevaluation step:")
	scope := object.NewScope(nil)
	builtins.Register(scope)
	runtime := evaluator.New(scope)
	_, err := runtime.Eval(root)
	if err != nil {
		err := errfmt.FormatEvaluationError(err, source, "<stdin>")
		fmt.Println(err)
		os.Exit(1)
	}

	// fmt.Println(obj.AsString())
}
