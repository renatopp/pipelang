package object

import (
	"fmt"

	"github.com/renatopp/langtools/utils"
	"github.com/renatopp/pipelang/internal/ast"
)

type Runner interface {
	RunCode(code []byte) (Object, error)
	RunFile(path string) (Object, error)
	RunAst(node ast.Node) (Object, error)
}

type Evaluator interface {
	RawEval(scope *Scope, node ast.Node) Object
	Call(scope *Scope, obj Object, args []Object) Object
	Operator(scope *Scope, op string, left, right Object) Object
}

type Scope struct {
	depth        int
	store        map[string]Object
	parent       *Scope
	stack        []ast.Node
	runner       Runner
	eval         Evaluator
	activeRecord ActiveRecord
}

func NewScope(r Runner) *Scope {
	return &Scope{
		depth:  0,
		store:  make(map[string]Object),
		parent: nil,
		runner: r,
		eval:   nil,
	}
}

func (s *Scope) New() *Scope {
	return &Scope{
		depth:  s.depth + 1,
		store:  make(map[string]Object),
		parent: s,
		runner: s.runner,
		eval:   s.eval,
	}
}

func (s *Scope) GetGlobal(name string) Object {
	if obj, ok := s.store[name]; ok {
		return obj
	}
	if s.parent != nil {
		return s.parent.GetGlobal(name)
	}
	return nil
}

func (s *Scope) GetLocal(name string) Object {
	if obj, ok := s.store[name]; ok {
		return obj
	}
	return nil
}

func (s *Scope) SetGlobal(name string, obj Object) {
	if _, ok := s.store[name]; ok {
		s.store[name] = obj
		return
	}

	if s.parent != nil {
		s.parent.SetGlobal(name, obj)
		return
	}
}

func (s *Scope) SetLocal(name string, obj Object) {
	s.store[name] = obj
}

func (s *Scope) PushNode(node ast.Node) {
	s.stack = append(s.stack, node)
}

func (s *Scope) PopNode() ast.Node {
	if len(s.stack) == 0 {
		return nil
	}

	node := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return node
}

// TODO: add stack trace and active record?
func (s *Scope) Interrupt(interruption *Interruption) Object {
	interruption.Stack = s.stack
	interruption.TriggeredScope = s
	return interruption
}

func (s *Scope) Keys() []string {
	keys := make([]string, 0, len(s.store))
	for key := range s.store {
		keys = append(keys, key)
	}
	return keys
}

func (s *Scope) Depth() int {
	return s.depth
}

func (s *Scope) Parent() *Scope {
	return s.parent
}

func (s *Scope) Runner() Runner {
	return s.runner
}

func (s *Scope) Eval() Evaluator {
	return s.eval
}

func (s *Scope) WithEval(eval Evaluator) *Scope {
	s.eval = eval
	return s
}

func (s *Scope) SetActiveRecord(ar ActiveRecord) {
	s.activeRecord = ar
}

func (s *Scope) ActiveRecord() ActiveRecord {
	return s.activeRecord
}

func (s *Scope) Print(name string) {
	parent := s
	i := 0
	println("--SCOPE ------------------------------------")
	println("- " + name)
	println("--------------------------------------------")
	for parent != nil {
		maxlen := 0
		for key := range parent.store {
			maxlen = max(maxlen, len(key))
		}

		for key, value := range parent.store {
			fmt.Printf("|%s: %s (%s)\n", utils.PadLeft(key, maxlen+1), value.AsString(), value.TypeId())
		}

		parent = parent.parent
		if parent != nil {
			println("|" + utils.PadCenter("~ parent ~", 40))
		}
		i++
	}
	println("--------------------------------------------")
}

func isRaise(obj Object) bool {
	intr, ok := obj.(*Interruption)
	return ok && intr.Category == RaiseId
}
