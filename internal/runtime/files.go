package runtime

import (
	"encoding/gob"
	"os"

	"github.com/renatopp/pipelang/internal/ast"

	"github.com/renatopp/pipelang/internal"
	"github.com/renatopp/pipelang/internal/errfmt"
	"github.com/renatopp/pipelang/internal/logs"
)

func init() {
	gob.Register(&CacheFile{})
}

type CacheFile struct {
	SourcePath string
	Ast        ast.Node
}

type SourceFile struct {
	hash       string
	sourcePath string
	cachePath  string
	source     []byte
	ast        ast.Node
}

func (f *SourceFile) Hash() string {
	return f.hash
}

func (f *SourceFile) SourcePath() string {
	return f.sourcePath
}

func (f *SourceFile) CachePath() string {
	return f.cachePath
}

func (f *SourceFile) LoadSource() ([]byte, error) {
	if f.source != nil {
		return f.source, nil
	}

	data, err := os.ReadFile(f.sourcePath)
	if err != nil {
		return nil, err
	}

	f.source = data
	return data, nil
}

func (f *SourceFile) LoadAst() (ast.Node, error) {
	if f.ast != nil {
		return f.ast, nil
	}

	source, err := f.LoadSource()
	if err != nil {
		return nil, err
	}

	logs.Print("[sourcefile] converting code to ast")
	// Pre-Lexing: extract all tokens as-is from the source code
	lexer := internal.NewPreLexer(source)
	tokens := lexer.All()
	if lexer.HasErrors() {
		logs.Print("[sourcefile] error pre-lexing (%v)", lexer.Errors())
		return nil, errfmt.FormatLexerErrors(lexer.Errors(), source, f.sourcePath)
	}

	// Post-Lexing: validate and optimize the tokens list
	postLexer := internal.NewPostLexer(tokens)
	if err := postLexer.Optimize(); err != nil {
		logs.Print("[sourcefile] error post-lexing (%v)", err)
		return nil, err
	}

	// Parsing: build the AST from the optimized tokens list
	parser := internal.NewPipeParser(postLexer)
	ast := parser.Parse()
	if parser.HasErrors() {
		logs.Print("[sourcefile] error parsing (%v)", parser.Errors())
		return nil, errfmt.FormatParserErrors(parser.Errors(), source, f.sourcePath)
	}

	// Post-Parsing for optimizations
	f.ast = internal.OptimizeAst(ast)

	return ast, nil
}
