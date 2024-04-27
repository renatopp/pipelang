package runtime

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/renatopp/pipelang/internal"
	"github.com/renatopp/pipelang/internal/ast"

	"github.com/renatopp/pipelang/internal/builtins"
	"github.com/renatopp/pipelang/internal/errfmt"
	"github.com/renatopp/pipelang/internal/evaluator"
	"github.com/renatopp/pipelang/internal/logs"
	o "github.com/renatopp/pipelang/internal/object"
)

// TODO: how make this concurrent safe?
type Runtime struct {
	globalScope *o.Scope
	objectCache map[string]o.Object
	fileCache   *FileCache
	loadStack   []string // file paths currently in execution, to detect circular imports
}

func New() *Runtime {
	r := &Runtime{}
	r.globalScope = o.NewScope(r)
	r.objectCache = make(map[string]o.Object)
	r.fileCache = NewFileCache()

	builtins.Register(r.globalScope)

	return r
}

func (r *Runtime) GlobalScope() *o.Scope {
	return r.globalScope
}

func (r *Runtime) LoadAst(code []byte) (ast.Node, error) {
	logs.Print("[runtime] running from code")

	file := &SourceFile{
		sourcePath: "<stdin>",
		source:     code,
	}

	ast, err := file.LoadAst()
	if err != nil {
		return nil, err
	}

	return ast, nil
}

func (r *Runtime) RunAst(node ast.Node) (o.Object, error) {
	logs.Print("[runtime] running from ast")

	obj, evalErr := r.evalAst(node)
	if evalErr != nil {
		// TODO: load source from node
		return nil, errfmt.FormatEvaluationError(evalErr, []byte{}, "<ast>")
	}

	return obj, nil
}

func (r *Runtime) RunCode(code []byte) (o.Object, error) {
	logs.Print("[runtime] running from code")

	file := &SourceFile{
		sourcePath: "<stdin>",
		source:     code,
	}

	ast, err := file.LoadAst()
	if err != nil {
		return nil, err
	}

	obj, evalErr := r.evalAst(ast)
	if evalErr != nil {
		// TODO: load source from node
		return nil, errfmt.FormatEvaluationError(evalErr, code, file.SourcePath())
	}

	return obj, nil
}

func (r *Runtime) RunFile(path string) (o.Object, error) {
	logs.Print("[runtime] running from file (%s)", path)
	path, err := r.getAbsolutePath(path)
	if err != nil {
		return nil, err
	}

	obj, ok := r.objectCache[path]
	if ok {
		return obj, nil
	}

	if err := r.pushFileStack(path); err != nil {
		return nil, err
	}
	defer r.popFileStack()

	file, err := r.fileCache.Load(path)
	if err != nil {
		return nil, err
	}

	obj, evalErr := r.evalAst(file.ast)
	if evalErr != nil {
		// TODO: load source from node
		source, err := file.LoadSource()
		if err != nil {
			source = []byte{}
		}

		return nil, errfmt.FormatEvaluationError(evalErr, source, path)
	}

	return obj, nil
}

func (r *Runtime) pushFileStack(path string) error {
	if i := slices.Index(r.loadStack, path); i >= 0 {
		return fmt.Errorf("circular import detected between '%s' and  '%s'", r.loadStack[i], path)
	}

	r.loadStack = append(r.loadStack, path)
	return nil
}

func (r *Runtime) popFileStack() {
	r.loadStack = slices.Delete(r.loadStack, len(r.loadStack)-1, len(r.loadStack))
}

func (r *Runtime) getAbsolutePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func (r *Runtime) evalAst(node ast.Node) (o.Object, *internal.Error) {
	eval := evaluator.New(r.globalScope)
	obj, err := eval.Eval(node)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
