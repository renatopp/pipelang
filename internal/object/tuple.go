package object

import (
	"fmt"
	"strings"
)

var TupleId = TypeIdentifier("Tuple")
var TupleTypeObj = NewTupleType()

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Number`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type TupleType struct {
	*BaseObjectType
}

func NewTupleType() *TupleType {
	return &TupleType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			TupleId,
		),
	}
}

func (o *TupleType) Instantiate(scope *Scope) Object {
	return scope.Interrupt(Raise("cannot instantiate type 'Tuple' manually"))
}

func (o *TupleType) Convert(scope *Scope, obj Object) Object {
	// TODO: convert map keys

	switch obj := obj.(type) {
	case *List:
		return NewTuple(obj.Elements...)

	case *Tuple:
		return obj
	}

	return NewTuple(obj)
}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type Tuple struct {
	*BaseObject
	Elements []Object
}

func NewTuple(e ...Object) *Tuple {
	return &Tuple{
		BaseObject: NewBaseObject(TupleTypeObj),
		Elements:   e,
	}
}

func (o *Tuple) AsBool() bool {
	return o.Elements[0].AsBool()
	// return true
}

func (o *Tuple) AsString() string {
	var elements []string
	for _, e := range o.Elements {
		elements = append(elements, e.AsRepr())
	}
	return fmt.Sprintf("(%s)", strings.Join(elements, ", "))
}

func (o *Tuple) AsInterface() any {
	return o.AsString()
}

func (o *Tuple) AsRepr() string {
	return o.AsString()
}
