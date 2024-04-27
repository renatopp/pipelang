package object

var ErrorId = TypeIdentifier("Error")
var ErrorTypeObj = NewErrorType()

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Number`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type ErrorType struct {
	*BaseObjectType
}

func NewErrorType() *ErrorType {
	t := &ErrorType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			ErrorId,
		),
	}

	t.AddMethod(Error_Msg)

	return t
}

var EmptyError = NewErrorFromString("Value is empty.")

func (o *ErrorType) Instantiate(scope *Scope) Object {
	// TODO: Implement
	return scope.Interrupt(Raise("cannot instantiate type 'Error' manually"))
}

func (o *ErrorType) Convert(scope *Scope, obj Object) Object {
	return NewError(obj)
}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type Error struct {
	*BaseObject
	Message string
}

func NewError(obj Object) *Error {
	return NewErrorFromString(obj.AsString())
}

func NewErrorFromString(msg string) *Error {
	e := &Error{
		BaseObject: NewBaseObject(ErrorTypeObj),
		Message:    msg,
	}

	return e
}

func (o *Error) AsBool() bool {
	return true
}

func (o *Error) AsString() string {
	return o.Message
}

func (o *Error) AsInterface() any {
	return o.AsString()
}

func (o *Error) AsRepr() string {
	return o.AsString()
}

// ----------------------------------------------------------------------------
// Instance Methods
// ----------------------------------------------------------------------------
var Error_Msg = NewBuiltinFunction("Msg", func(scope *Scope, args ...Object) Object {
	this := args[0].(*Error)
	return NewString(this.Message)
})
