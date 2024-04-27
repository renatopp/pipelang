package object

import (
	"fmt"

	"github.com/renatopp/pipelang/internal/ast"
)

var FunctionId = TypeIdentifier("Function")
var FunctionTypeObj = NewFunctionType()

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Function`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type FunctionType struct {
	*BaseObjectType
}

func NewFunctionType() *FunctionType {
	return &FunctionType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			FunctionId,
		),
	}
}

func (o *FunctionType) Instantiate(scope *Scope) Object {
	return Zero
}

func (o *FunctionType) Convert(scope *Scope, obj Object) Object {
	return scope.Interrupt(Raise("type 'Function' does not support conversion"))
}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type Function struct {
	*BaseObject
	Name       string
	Parameters []ast.Node
	Body       ast.Node
	Scope      *Scope
}

func NewFunction(name string, params []ast.Node, body ast.Node, scope *Scope) *Function {
	return &Function{
		BaseObject: NewBaseObject(FunctionTypeObj),
		Name:       name,
		Parameters: params,
		Body:       body,
		Scope:      scope,
	}
}

func (o *Function) AsBool() bool {
	return false
}

func (o *Function) AsString() string {
	if o.Name != "" {
		return "<Function:" + o.Name + ">"
	}
	return "<Function>"
}

func (o *Function) AsRepr() string {
	return o.AsString()
}

// ----------------------------------------------------------------------------
// Builtin Functions
// ----------------------------------------------------------------------------
type BuiltinFn func(scope *Scope, args ...Object) Object
type BuiltinFunction struct {
	*BaseObject
	Name      string
	Fn        BuiltinFn
	Signature []*Param
}

func NewBuiltinFunction(name string, fn BuiltinFn) *BuiltinFunction {
	return &BuiltinFunction{
		BaseObject: NewBaseObject(FunctionTypeObj),
		Name:       name,
		Fn:         fn,
		Signature:  []*Param{},
	}
}

// Shortcut for builtin functions with parameters
func F(fn BuiltinFn, name, docs string, params ...*Param) *BuiltinFunction {
	b := NewBuiltinFunction(name, fn)
	b.SetDocs(docs)
	b.WithSignature(params...)
	return b
}

func (o *BuiltinFunction) AsBool() bool {
	return true
}

func (o *BuiltinFunction) AsString() string {
	return fmt.Sprintf("<go:Function:%s>", o.Name)
}

func (o *BuiltinFunction) AsInterface() any {
	return o.Fn
}

func (o *BuiltinFunction) AsRepr() string {
	return o.AsString()
}

func (o *BuiltinFunction) Call(scope *Scope, args ...Object) Object {
	if err := o.Check(scope, args...); err != nil {
		return scope.Interrupt(RaiseWith(err))
	}

	return o.Fn(scope, args...)
}

func (o *BuiltinFunction) Check(scope *Scope, args ...Object) Object {
	j := 0
	for i, p := range o.Signature {

		to := j + 1
		if p.Spread {
			remaining := len(o.Signature) - i - 1
			to = len(args) - i - remaining
		}

		for ; j < to; j++ {
			var arg Object
			if j < len(args) {
				arg = args[j]
			}

			for _, v := range p.Validations {
				if err := v(p, arg); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (o *BuiltinFunction) WithDocs(d string) *BuiltinFunction {
	o.SetDocs(d)
	return o
}

func (o *BuiltinFunction) WithSignature(params ...*Param) *BuiltinFunction {
	o.Signature = params
	return o
}

type Param struct {
	Name        string
	Spread      bool
	Validations []ValidationFn
}

func P(name string, validations ...ValidationFn) *Param {
	return &Param{
		Name:        name,
		Validations: validations,
	}
}

func (p *Param) AsSpread() *Param {
	p.Spread = true
	return p
}
