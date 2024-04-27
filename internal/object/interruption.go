package object

// INTERNAL USE ONLY

import (
	"fmt"

	"github.com/renatopp/pipelang/internal/ast"
)

var InterruptionId TypeIdentifier = "interruption"

type InterruptionType string

var (
	ReturnId           InterruptionType = "return"
	RaiseId            InterruptionType = "raise"
	YieldId            InterruptionType = "yield"
	BreakId            InterruptionType = "break"
	ContinueId         InterruptionType = "continue"
	DeferredFunctionId InterruptionType = "deferred_function"
)

type Interruption struct {
	*BaseObject
	Category       InterruptionType
	TriggeredScope *Scope
	Value          Object
	Context        any
	Stack          []ast.Node
}

func (i *Interruption) Id() string                 { return "" }
func (i *Interruption) TypeId() TypeIdentifier     { return InterruptionId }
func (i *Interruption) Type() ObjectType           { return i.Value.Type() }
func (i *Interruption) SetDocs(string)             {}
func (i *Interruption) Docs() string               { return "" }
func (i *Interruption) SetProperty(string, Object) {}
func (i *Interruption) GetProperty(string) Object  { return nil }
func (i *Interruption) AsBool() bool               { return i.Value.AsBool() }
func (i *Interruption) AsString() string           { return string(i.Category) + ": " + i.Value.AsString() }

func Return(msg string, v ...any) *Interruption {
	return ReturnWith(NewString(fmt.Sprintf(msg, v...)))
}
func ReturnWith(value Object) *Interruption {
	return &Interruption{
		Category: ReturnId,
		Value:    value,
	}
}

func Raise(msg string, v ...any) *Interruption {
	m := NewString(fmt.Sprintf(msg, v...))
	return RaiseWith(m)
}
func RaiseWith(value Object) *Interruption {
	return &Interruption{
		Category: RaiseId,
		Value:    NewError(value),
	}
}

func Yield(msg string, v ...any) *Interruption {
	return YieldWith(NewString(fmt.Sprintf(msg, v...)))
}
func YieldWith(value Object) *Interruption {
	return &Interruption{
		Category: YieldId,
		Value:    value,
	}
}

// func Break() *Interruption {
// 	return &Interruption{
// 		Category: BreakId,
// 	}
// }

// func Continue() *Interruption {
// 	return &Interruption{
// 		Category: ContinueId,
// 	}
// }
