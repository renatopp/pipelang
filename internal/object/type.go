package object

var TypeId = TypeIdentifier("Type")
var TypeTypeObj = NewType()

type Type struct {
	*BaseObject
}

func NewType() *Type {
	tp := &Type{
		BaseObject: NewBaseObject(nil),
	}
	tp.otype = tp
	return tp
}

func (o *Type) TypeId() TypeIdentifier {
	return TypeId
}

func (o *Type) Instantiate(scope *Scope) Object {
	// TODO: Throw
	return nil
}

func (o *Type) Convert(scope *Scope, obj Object) Object {
	// TODO: Throw
	return nil
}

func (o *Type) AsBool() bool {
	return true
}

func (o *Type) AsString() string {
	return string(TypeId)
}

func (o *Type) AsRepr() string {
	return o.AsString()
}
