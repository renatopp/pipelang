package object

import (
	"fmt"
)

var MaybeId = TypeIdentifier("Maybe")
var MaybeTypeObj = NewMaybeType()

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Number`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type MaybeType struct {
	*BaseObjectType
}

func NewMaybeType() *MaybeType {
	t := &MaybeType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			MaybeId,
		),
	}

	t.AddMethod(Maybe_Ok)
	t.AddMethod(Maybe_Value)
	t.AddMethod(Maybe_Error)
	t.AddMethod(Maybe_Result)

	return t
}

func (o *MaybeType) Instantiate(scope *Scope) Object {
	return scope.Interrupt(Raise("cannot instantiate type 'Maybe' manually"))
}

func (o *MaybeType) Convert(scope *Scope, obj Object) Object {
	return NewMaybeWithType(obj, obj.TypeId())
}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type Maybe struct {
	*BaseObject
	Ok        bool
	Value     Object
	Error     Object
	ValueType TypeIdentifier
}

func NewMaybeWithType(obj Object, tp TypeIdentifier) *Maybe {
	if obj.TypeId() == MaybeId {
		obj = obj.(*Maybe).Result()
	}

	m := &Maybe{
		BaseObject: NewBaseObject(MaybeTypeObj),
		ValueType:  tp,
	}
	m.Set(obj)

	return m
}

func NewMaybe(obj Object) *Maybe {
	if obj.TypeId() == InterruptionId {
		obj = obj.(*Interruption).Value
	}

	return NewMaybeWithType(obj, obj.TypeId())
}

func (o *Maybe) AsBool() bool {
	return o.Ok
}

func (o *Maybe) AsString() string {
	ok := "ok"
	if !o.Ok {
		ok = "error"
	}

	return fmt.Sprintf("Maybe(%s, %s)", o.ValueType, ok)
}

func (o *Maybe) AsInterface() any {
	return o.AsString()
}

func (o *Maybe) IsType(t Object) bool {
	return t.TypeId() == o.ValueType || t.TypeId() == ErrorId
}

func (o *Maybe) Set(obj Object) {
	var val Object = False
	var err Object = False

	if obj.TypeId() == ErrorId {
		err = obj
	} else {
		val = obj
		o.ValueType = obj.TypeId()
	}

	o.Ok = val != False
	o.Value = val
	o.Error = err
}

func (o *Maybe) Result() Object {
	if o.Ok {
		return o.Value
	}
	return o.Error
}

func (o *Maybe) AsRepr() string {
	return o.AsString()
}

// ----------------------------------------------------------------------------
// Instance Methods
// ----------------------------------------------------------------------------
var Maybe_Ok = NewBuiltinFunction("Ok", func(scope *Scope, args ...Object) Object {
	this := args[0].(*Maybe)
	return NewBoolean(this.Ok)
})

var Maybe_Value = NewBuiltinFunction("Value", func(scope *Scope, args ...Object) Object {
	this := args[0].(*Maybe)
	if this.Ok {
		return this.Value
	}
	return False
})

var Maybe_Error = NewBuiltinFunction("Error", func(scope *Scope, args ...Object) Object {
	this := args[0].(*Maybe)
	if !this.Ok {
		return this.Error
	}
	return False
})

var Maybe_Result = NewBuiltinFunction("Result", func(scope *Scope, args ...Object) Object {
	this := args[0].(*Maybe)
	return this.Result()
})
