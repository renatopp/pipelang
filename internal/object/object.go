package object

import (
	"github.com/renatopp/pipelang/internal/random"
)

type TypeIdentifier string

// Represent any instance in Pipe
type Object interface {
	Id() string             // Unique identifier for the instance
	TypeId() TypeIdentifier // Type ID, used to identify which struct we should use
	Type() ObjectType

	Docs() string
	SetDocs(string)

	Parent() Object
	SetParent(Object)

	GetProperty(string) Object
	SetProperty(string, Object)

	OnIndex(*Scope, *Tuple) Object
	OnIndexAssign(*Scope, *Tuple, Object) Object
	OnOperator(*Scope, string, Object) Object

	AsBool() bool
	AsString() string
	AsRepr() string
	AsInterface() any
}

// Represent type instances, like `Number`, `String` or even `Type`
type ObjectType interface {
	Object

	AddMethod(*BuiltinFunction)
	Convert(*Scope, Object) Object
	Instantiate(*Scope) Object
}

// ----------------------------------------------------------------------------
// BaseObject Definition
// ----------------------------------------------------------------------------
type BaseObject struct {
	id         string
	docs       string
	otype      ObjectType
	properties map[string]Object
	parent     Object // Used for access, like `a.b.c` where `a` is the parent of `b`
}

func NewBaseObject(otype ObjectType) *BaseObject {
	return &BaseObject{
		id:         random.NanoId(10),
		otype:      otype,
		properties: make(map[string]Object),
	}
}

func (o *BaseObject) Id() string {
	return o.id
}

func (o *BaseObject) TypeId() TypeIdentifier {
	if o.otype == nil {
		return ""
	}
	return o.otype.TypeId()
}

func (o *BaseObject) Type() ObjectType {
	return o.otype
}

func (o *BaseObject) Docs() string {
	return o.docs
}

func (o *BaseObject) SetDocs(docs string) {
	o.docs = docs
}

func (o *BaseObject) Parent() Object {
	return o.parent
}

func (o *BaseObject) SetParent(p Object) {
	o.parent = p
}

func (o *BaseObject) OnIndex(scope *Scope, t *Tuple) Object {
	return scope.Interrupt(Raise("type '%s' does not support indexing", o.TypeId()))
}

func (o *BaseObject) OnIndexAssign(scope *Scope, t *Tuple, value Object) Object {
	return scope.Interrupt(Raise("type '%s' does not support indexing assignment", o.TypeId()))
}

func (o *BaseObject) OnOperator(scope *Scope, op string, right Object) Object {
	return scope.Interrupt(Raise("type '%s' does not support operator '%s'", o.TypeId(), op))
}

func (o *BaseObject) GetProperty(name string) Object {
	prop, ok := o.properties[name]
	if ok {
		return prop
	}

	parent := o.Type()
	myId := o.Id()
	parentId := parent.Id()
	if myId == parentId {
		return nil
	}

	return parent.GetProperty(name)
}

func (o *BaseObject) SetProperty(name string, value Object) {
	o.properties[name] = value
}

func (o *BaseObject) AsBool() bool {
	return false
}

func (o *BaseObject) AsString() string {
	return ""
}

func (o *BaseObject) AsRepr() string {
	return o.AsString()
}

func (o *BaseObject) AsInterface() any {
	return o.AsString()
}

// ----------------------------------------------------------------------------
// BaseObjectType Definition
// ----------------------------------------------------------------------------
type BaseObjectType struct {
	*BaseObject

	typeId TypeIdentifier
}

func NewBaseObjectType(baseObject *BaseObject, typeId TypeIdentifier) *BaseObjectType {
	return &BaseObjectType{
		BaseObject: baseObject,
		typeId:     typeId,
	}
}

func (o *BaseObjectType) TypeId() TypeIdentifier {
	return o.typeId
}

func (o *BaseObjectType) AsBool() bool {
	return true
}

func (o *BaseObjectType) AsString() string {
	return string(o.typeId)
}

func (o *BaseObjectType) AsRepr() string {
	return o.AsString()
}

func (o *BaseObject) AddMethod(fn *BuiltinFunction) {
	o.properties[fn.Name] = fn
}

// ---
