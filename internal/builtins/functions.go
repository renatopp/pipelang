package builtins

import (
	o "github.com/renatopp/pipelang/internal/object"
)

func RegisterBuiltinFunctions(s *o.Scope) {
	setFunction(s, o.Printf)
	setFunction(s, o.Printfln)
	setFunction(s, o.Sprintf)
	setFunction(s, o.Sprintfln)
	setFunction(s, o.Print)
	setFunction(s, o.Println)
	setFunction(s, o.Sprint)
	setFunction(s, o.Sprintln)
	setFunction(s, o.Range)
	setFunction(s, o.Filter)
	setFunction(s, o.Each)
	setFunction(s, o.Map)
	setFunction(s, o.Reduce)
	setFunction(s, o.Sum)
	setFunction(s, o.SumBy)
	setFunction(s, o.Count)
	setFunction(s, o.CountBy)
	setFunction(s, Import)
}

var Import = o.NewBuiltinFunction("import", func(s *o.Scope, args ...o.Object) o.Object {
	path := args[0].AsString()
	obj, err := s.Runner().RunFile(path)
	if err != nil {
		return s.Interrupt(o.Raise(err.Error()))
	}

	return obj
})

var Empty = o.NewBuiltinFunction("empty", func(s *o.Scope, args ...o.Object) o.Object {
	return o.NewMaybeWithType(o.EmptyError, args[0].TypeId())
})
