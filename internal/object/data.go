package object

import (
	"fmt"

	"github.com/renatopp/pipelang/internal/ast"
)

var DataId = TypeIdentifier("Data")

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Data`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
// The data type should be created by user
type DataType struct {
	*BaseObjectType
	Name       string
	Attributes map[string]ast.Node
	Methods    map[string]*Function
}

func NewDataType(name string, attributes map[string]ast.Node, methods map[string]*Function) *DataType {
	dt := &DataType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			DataId,
		),
		Name:       name,
		Attributes: attributes,
		Methods:    methods,
	}

	for name, method := range methods {
		dt.SetProperty(name, method)
	}

	return dt
}

func (o *DataType) AsString() string {
	if o.Name == "" {
		return "<Data>"
	}

	return fmt.Sprintf("<Data:%s>", o.Name)
}

func (o *DataType) Instantiate(scope *Scope) Object {
	return &Data{
		BaseObject: NewBaseObject(o),
	}
}

func (o *DataType) Convert(scope *Scope, obj Object) Object {
	return scope.Interrupt(Raise("type 'Data' does not support conversion"))
}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type Data struct {
	*BaseObject
}

func (o *Data) AsBool() bool {
	return true
}

func (o *Data) AsString() string {
	tp := o.Type().(*DataType)

	if tp.Name == "" {
		return "<Data Instance>"
	} else {
		return "<Data Instance:" + o.Type().(*DataType).Name + ">"
	}
}

func (o *Data) AsRepr() string {
	return o.AsString()
}
