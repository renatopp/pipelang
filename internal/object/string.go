package object

import (
	"fmt"
	"strings"
)

var StringId = TypeIdentifier("String")
var StringTypeObj = NewStringType()
var EmptyString = NewString("")

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `String`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type StringType struct {
	*BaseObjectType
}

func NewStringType() *StringType {
	return &StringType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			StringId,
		),
	}
}

func (o *StringType) Instantiate(scope *Scope) Object {
	return EmptyString
}

func (o *StringType) Convert(scope *Scope, obj Object) Object {
	switch obj := obj.(type) {
	case *Stream:
		result := ""
		ret := obj.Resolve(func(obj Object) Object {
			result += obj.AsString()
			return nil
		})
		if isRaise(ret) {
			return ret
		}
		return NewString(result)

	default:
		return NewString(obj.AsString())
	}

}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type String struct {
	*BaseObject
	Value string
}

func NewString(value string) *String {
	return &String{
		BaseObject: NewBaseObject(StringTypeObj),
		Value:      value,
	}
}

func (o *String) OnIndex(scope *Scope, t *Tuple) Object {
	if len(t.Elements) == 1 {
		return String_Get.Call(scope, o, t.Elements[0])
	} else if len(t.Elements) == 2 {
		return String_Sub.Call(scope, o, t.Elements[0], t.Elements[1])
	}
	return scope.Interrupt(Raise("invalid string index '%s'", t.AsString()))
}

func (o *String) Copy() Object {
	return NewString(o.Value)
}

func (o *String) OnOperator(scope *Scope, op string, right Object) Object {
	a := o.Value
	b := right.(*String).Value

	switch op {
	case "+":
		return NewString(a + b)
	case "==":
		return NewBoolean(a == b)
	case "!=":
		return NewBoolean(a != b)
	case ">":
		return NewBoolean(a > b)
	case "<":
		return NewBoolean(a < b)
	case ">=":
		return NewBoolean(a >= b)
	case "<=":
		return NewBoolean(a <= b)
	case "<=>":
		if a < b {
			return MinusOne
		} else if a > b {
			return One
		}
		return Zero
	default:
		return scope.Interrupt(Raise("type 'String' does not support operator '%s'", op))
	}
}

func (o *String) AsBool() bool {
	return o.Value != ""
}

func (o *String) AsString() string {
	return o.Value
}

func (o *String) AsInterface() any {
	return o.Value
}

func (o *String) AsRepr() string {
	v := o.Value
	v = strings.ReplaceAll(v, "\n", "\\n")
	v = strings.ReplaceAll(v, "'", "\\'")
	return fmt.Sprintf("'%s'", v)
}
