package builtins

import (
	o "github.com/renatopp/pipelang/internal/object"
)

func Register(s *o.Scope) {
	RegisterBuiltinTypes(s)
	RegisterBuiltinFunctions(s)
	RegisterBuiltinModules(s)
}

func setFunction(s *o.Scope, fn *o.BuiltinFunction) {
	s.SetLocal(fn.Name, fn)
}
