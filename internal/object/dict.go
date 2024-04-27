package object

import (
	"fmt"
	"strings"
)

var DictId = TypeIdentifier("Dict")
var DictTypeObj = NewDictType()

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Number`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type DictType struct {
	*BaseObjectType
}

func NewDictType() *DictType {
	return &DictType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			DictId,
		),
	}
}

func (o *DictType) Instantiate(scope *Scope) Object {
	return scope.Interrupt(Raise("cannot instantiate type 'Dict' manually"))
}

func (o *DictType) Convert(scope *Scope, obj Object) Object {
	// TODO: convert map keys

	// switch obj := obj.(type) {
	// case *Dict:
	// 	return NewDict(obj.Elements...)

	// case *Tuple:
	// 	return NewDict(obj.Elements...)
	// }

	return NewDictFromList(obj)
}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type Dict struct {
	*BaseObject
	Elements map[string]Object
}

func NewDict(elements map[string]Object) *Dict {
	return &Dict{
		BaseObject: NewBaseObject(DictTypeObj),
		Elements:   elements,
	}
}

// [Key, Value, Key, Value, ...]
func NewDictFromList(e ...Object) *Dict {
	d := &Dict{
		BaseObject: NewBaseObject(DictTypeObj),
		Elements:   map[string]Object{},
	}

	for i := 0; i < len(e); i += 2 {
		key := e[i].AsString()
		value := e[i+1]
		d.Elements[key] = value
	}

	return d
}

func (o *Dict) OnIndex(scope *Scope, t *Tuple) Object {
	if len(t.Elements) == 1 {
		return Dict_Get.Call(scope, o, t.Elements[0])
	}
	return scope.Interrupt(Raise("invalid string index '%s'", t.AsString()))
}

func (o *Dict) OnIndexAssign(scope *Scope, t *Tuple, value Object) Object {
	if len(t.Elements) == 1 {
		return Dict_Set.Call(scope, o, t.Elements[0], value)
	}
	return scope.Interrupt(Raise("invalid string index '%s'", t.AsString()))
}

func (o *Dict) AsBool() bool {
	return true
}

func (o *Dict) AsString() string {
	var elements []string
	for k, v := range o.Elements {
		elements = append(elements, fmt.Sprintf("%s=%s", k, v.AsString()))
	}
	return fmt.Sprintf("{%s}", strings.Join(elements, ", "))
}

func (o *Dict) AsInterface() any {
	return o.AsString()
}
