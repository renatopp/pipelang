package builtins

import (
	o "github.com/renatopp/pipelang/internal/object"
)

// import "reflect"

func RegisterBuiltinTypes(s *o.Scope) {
	s.SetLocal("Type", o.TypeTypeObj)
	s.SetLocal("Number", o.NumberTypeObj)
	s.SetLocal("String", o.StringTypeObj)
	s.SetLocal("Boolean", o.BooleanTypeObj)
	s.SetLocal("Function", o.FunctionTypeObj)
	s.SetLocal("Tuple", o.TupleTypeObj)
	s.SetLocal("List", o.ListTypeObj)
	s.SetLocal("Maybe", o.MaybeTypeObj)
	s.SetLocal("Error", o.ErrorTypeObj)
	s.SetLocal("Stream", o.ErrorTypeObj)
}
