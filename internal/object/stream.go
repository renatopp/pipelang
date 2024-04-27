package object

var StreamId = TypeIdentifier("Stream")
var StreamTypeObj = NewStreamType()

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Number`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type StreamType struct {
	*BaseObjectType
}

func NewStreamType() *StreamType {
	t := &StreamType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			StreamId,
		),
	}

	t.AddMethod(Stream_Next)
	t.AddMethod(Stream_Finished)

	return t
}

func (o *StreamType) Instantiate(scope *Scope) Object {
	return scope.Interrupt(Raise("cannot instantiate type 'Stream' manually"))
}

func (o *StreamType) Convert(scope *Scope, obj Object) Object {
	switch obj.TypeId() {
	case StreamId:
		return obj

	case StringId:
		return String_Chars.Call(scope, obj)

	case ListId:
		return List_Elements.Call(scope, obj)

	default:
		return NewInternalStream(func(s *Scope) Object {
			return YieldWith(obj)
		}, scope)
	}

}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type Stream struct {
	*BaseObject
	Finished   bool
	Scope      *Scope
	Fn         *Function
	InternalFn func(*Scope) Object
	iteration  *StreamIteration // used only to indicate the evaluator to call the stream
}

// Scope is the fixed scope of the function
func NewStream(fn *Function, scope *Scope) *Stream {
	s := &Stream{
		BaseObject: NewBaseObject(StreamTypeObj),
		Finished:   false,
		Scope:      scope,
		Fn:         fn,
	}
	s.iteration = &StreamIteration{
		Stream: s,
	}

	return s
}
func NewInternalStream(fn func(*Scope) Object, scope *Scope) *Stream {
	s := &Stream{
		BaseObject: NewBaseObject(StreamTypeObj),
		Finished:   false,
		Scope:      scope,
		InternalFn: fn,
	}
	s.iteration = &StreamIteration{
		Stream: s,
	}

	return s
}

func (o *Stream) Resolve(fn func(Object) Object) Object {
	for {
		maybe := Stream_Next.Call(o.Scope, o)
		if isRaise(maybe) {
			return maybe
		}
		if o.Finished {
			return nil
		}

		value := maybe.(*Maybe).Value
		ret := fn(value)
		if isRaise(ret) {
			return ret
		}
	}
}

func (o *Stream) AsBool() bool {
	return true
}

func (o *Stream) AsString() string {
	return "<stream>"
}

func (o *Stream) AsRepr() string {
	return o.AsString()
}

func (o *Stream) AsInterface() any {
	return o.AsString()
}

// ----------------------------------------------------------------------------
// Instance Methods
// ----------------------------------------------------------------------------
var Stream_Next = NewBuiltinFunction("Next", func(scope *Scope, args ...Object) Object {
	this := args[0].(*Stream)
	if this.Finished {
		return NewMaybe(NewErrorFromString("Stream finished"))
	}

	var ret Object
	if this.Fn != nil {
		ret = scope.Eval().RawEval(this.Scope, this.Fn.Body)
	} else {
		ret = this.InternalFn(this.Scope)
	}

	if isRaise(ret) {
		return ret
	}

	if t := asYield(ret); t != nil {
		return NewMaybe(t.Value)
	}

	this.Finished = true
	return NewMaybe(NewErrorFromString("Stream finished"))
})

var Stream_Finished = NewBuiltinFunction("Finished", func(scope *Scope, args ...Object) Object {
	this := args[0].(*Stream)
	return NewBoolean(this.Finished)
})

// func streamFinishedError() o.Object {
// 	return o.NewMaybe(o.NewErrorFromString("Stream finished"))
// }

func asYield(obj Object) *Interruption {
	intr, ok := obj.(*Interruption)
	if ok && intr.Category == YieldId {
		return intr
	}
	return nil
}
