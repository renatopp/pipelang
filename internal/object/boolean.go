package object

import (
	"fmt"
)

var BooleanId = TypeIdentifier("Boolean")
var BooleanTypeObj = NewBooleanType()
var False = internalBoolean(false)
var True = internalBoolean(true)

func GetBoolean(value bool) Object {
	if value {
		return True
	}
	return False
}

func GetInverseBoolean(value bool) Object {
	if value {
		return False
	}
	return True
}

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Boolean`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type BooleanType struct {
	*BaseObjectType
}

func NewBooleanType() *BooleanType {
	return &BooleanType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			BooleanId,
		),
	}
}

func (o *BooleanType) Instantiate(scope *Scope) Object {
	return False
}

func (o *BooleanType) Convert(scope *Scope, obj Object) Object {
	return NewBoolean(obj.AsBool())
}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type Boolean struct {
	*BaseObject
	Value bool
}

func internalBoolean(value bool) *Boolean {
	return &Boolean{
		BaseObject: NewBaseObject(BooleanTypeObj),
		Value:      value,
	}
}

func NewBoolean(value bool) Object {
	if value {
		return True
	}
	return False
}

func (o *Boolean) OnOperator(scope *Scope, op string, right Object) Object {
	a := o.Value
	b := right.(*Boolean).Value

	switch op {
	case "==":
		return NewBoolean(a == b)
	case "!=":
		return NewBoolean(a != b)
	default:
		return scope.Interrupt(Raise("type 'Boolean' does not support operator '%s'", op))
	}
}

func (o *Boolean) AsBool() bool {
	return o.Value
}

func (o *Boolean) AsString() string {
	return fmt.Sprintf("%t", o.Value)
}

func (o *Boolean) AsRepr() string {
	return o.AsString()
}

func (o *Boolean) AsInterface() any {
	return o.Value
}
