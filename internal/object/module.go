package object

import "fmt"

var ModuleId = TypeIdentifier("Module")

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Data`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type ModuleType struct {
	*BaseObjectType
	Name string
}

func NewModuleType(name string) *ModuleType {
	return &ModuleType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			ModuleId,
		),
		Name: name,
	}
}

func (o *ModuleType) AsString() string {
	if o.Name == "" {
		return "<Module>"
	}

	return fmt.Sprintf("<Module:%s>", o.Name)
}

func (o *ModuleType) Instantiate(scope *Scope) Object {
	return scope.Interrupt(Raise("cannot instantiate type 'Error' manually"))
}

func (o *ModuleType) Convert(scope *Scope, obj Object) Object {
	return scope.Interrupt(Raise("type 'Module' does not support conversion"))
}
