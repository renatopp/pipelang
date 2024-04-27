package builtins

import o "github.com/renatopp/pipelang/internal/object"

func RegisterBuiltinModules(s *o.Scope) {
	addModule(s, o.Module_Math)
}

func addModule(s *o.Scope, module *o.ModuleType) {
	s.SetLocal(module.Name, module)
}
