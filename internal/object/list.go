package object

import (
	"fmt"
	"strings"
)

var ListId = TypeIdentifier("List")
var ListTypeObj = NewListType()

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Number`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type ListType struct {
	*BaseObjectType
}

func NewListType() *ListType {
	return &ListType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			ListId,
		),
	}
}

func (o *ListType) Instantiate(scope *Scope) Object {
	return scope.Interrupt(Raise("cannot instantiate type 'List' manually"))
}

func (o *ListType) Convert(scope *Scope, obj Object) Object {
	// TODO: convert map keys

	switch obj := obj.(type) {
	case *List:
		return NewList(obj.Elements...)

	case *Tuple:
		return NewList(obj.Elements...)

	case *Stream:
		result := []Object{}
		ret := obj.Resolve(func(obj Object) Object {
			switch obj := obj.(type) {
			case *Tuple:
				result = append(result, obj.Elements[0])
			default:
				result = append(result, obj)
			}
			return nil
		})
		if isRaise(ret) {
			return ret
		}
		return NewList(result...)
	}

	return NewList(obj)
}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type List struct {
	*BaseObject
	Elements []Object
}

func NewList(e ...Object) *List {
	return &List{
		BaseObject: NewBaseObject(ListTypeObj),
		Elements:   e,
	}
}

func (o *List) OnIndex(scope *Scope, t *Tuple) Object {
	if len(t.Elements) == 1 {
		return List_Get.Call(scope, o, t.Elements[0])
	} else if len(t.Elements) == 2 {
		return List_Sub.Call(scope, o, t.Elements[0], t.Elements[1])
	}
	return scope.Interrupt(Raise("invalid string index '%s'", t.AsString()))
}

func (o *List) OnIndexAssign(scope *Scope, t *Tuple, value Object) Object {
	if len(t.Elements) == 1 {
		return List_Set.Call(scope, o, t.Elements[0], value)
	}

	return scope.Interrupt(Raise("invalid string index '%s'", t.AsString()))
}

func (o *List) AsBool() bool {
	return true
}

func (o *List) AsString() string {
	var elements []string
	for _, e := range o.Elements {
		elements = append(elements, e.AsRepr())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (o *List) AsInterface() any {
	return o.AsString()
}

func (o *List) AsRepr() string {
	return o.AsString()
}
