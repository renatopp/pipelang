package object

// INTERNAL USE ONLY

var StreamIterationId = TypeIdentifier("stream_iteration")

type StreamIteration struct {
	*BaseObject
	Stream *Stream
}

func (i *StreamIteration) Id() string                 { return "" }
func (i *StreamIteration) TypeId() TypeIdentifier     { return StreamIterationId }
func (i *StreamIteration) Type() ObjectType           { return i.Stream.Type() }
func (i *StreamIteration) SetDocs(string)             {}
func (i *StreamIteration) Docs() string               { return "" }
func (i *StreamIteration) SetProperty(string, Object) {}
func (i *StreamIteration) GetProperty(string) Object  { return nil }
func (i *StreamIteration) AsBool() bool               { return i.Stream.AsBool() }
func (i *StreamIteration) AsString() string           { return i.Stream.AsString() }
