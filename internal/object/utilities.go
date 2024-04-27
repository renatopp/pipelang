package object

import "fmt"

type ValidationFn func(p *Param, arg Object) Object
type Validator struct{}

var V = &Validator{}

func (v *Validator) Type(typeId TypeIdentifier) ValidationFn {
	return func(p *Param, arg Object) Object {
		if arg == nil {
			return NewErrorFromString(fmt.Sprintf("Expected argument '%s'.", p.Name))
		}

		if arg.TypeId() != typeId {
			return NewErrorFromString(fmt.Sprintf("Expected argument '%s' to be of type '%s', received '%s' instead.", p.Name, typeId, arg.TypeId()))
		}

		return nil
	}
}
