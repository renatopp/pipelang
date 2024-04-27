package object

import (
	"fmt"
	"math"
	"unicode/utf8"
)

var NumberId = TypeIdentifier("Number")
var NumberTypeObj = NewNumberType()
var Zero = NewNumber(0)
var One = NewNumber(1)
var Two = NewNumber(2)
var MinusOne = NewNumber(-1)

// ----------------------------------------------------------------------------
// Type Definition - represents the type instance in Pipe, like `Number`,
// `String` or even `Type`.
// ----------------------------------------------------------------------------
type NumberType struct {
	*BaseObjectType
}

func NewNumberType() *NumberType {
	return &NumberType{
		BaseObjectType: NewBaseObjectType(
			NewBaseObject(TypeTypeObj),
			NumberId,
		),
	}
}

func (o *NumberType) Instantiate(scope *Scope) Object {
	return Zero
}

func (o *NumberType) Convert(scope *Scope, obj Object) Object {
	switch obj := obj.(type) {
	case *Number:
		return obj

	case *String:
		if obj.Value == "" {
			return Zero
		}
		firstRune, _ := utf8.DecodeRuneInString(obj.Value)
		return NewNumber(float64(firstRune))

	case *List:
		if len(obj.Elements) == 0 {
			return Zero
		}
		return o.Convert(scope, obj.Elements[0])

	case *Stream:
		ret := obj.Resolve(func(obj Object) Object {
			switch obj := obj.(type) {
			case *Tuple:
				return obj.Elements[0]
			default:
				return obj
			}
		})
		if isRaise(ret) {
			return ret
		}
		return o.Convert(scope, ret)

	default:
		return Zero
	}
}

// ----------------------------------------------------------------------------
// Instance Definition - represents the instance of a particular type in Pipe,
// like `1` and `'foo'`.
// ----------------------------------------------------------------------------
type Number struct {
	*BaseObject
	Value float64
}

func NewNumber(value float64) *Number {
	return &Number{
		BaseObject: NewBaseObject(NumberTypeObj),
		Value:      value,
	}
}

func (o *Number) OnOperator(scope *Scope, op string, right Object) Object {
	a := o.Value
	b := right.(*Number).Value

	switch op {
	case "+":
		return NewNumber(a + b)
	case "-":
		return NewNumber(a - b)
	case "*":
		return NewNumber(a * b)
	case "/":
		return NewNumber(a / b)
	case "%":
		return NewNumber(float64(int(a) % int(b)))
	case "^":
		return NewNumber(math.Pow(a, b))
	case "==":
		return NewBoolean(a == b)
	case "!=":
		return NewBoolean(a != b)
	case "<":
		return NewBoolean(a < b)
	case ">":
		return NewBoolean(a > b)
	case "<=":
		return NewBoolean(a <= b)
	case ">=":
		return NewBoolean(a >= b)
	case "<=>":
		if a < b {
			return MinusOne
		} else if a > b {
			return One
		}
		return Zero
	default:
		return scope.Interrupt(Raise("type 'Number' does not support operator '%s'", op))
	}
}

func (o *Number) AsBool() bool {
	return o.Value != 0
}

func (o *Number) AsString() string {
	v := o.Value
	if math.IsInf(v, 1) {
		return "inf"
	} else if math.IsInf(v, -1) {
		return "-inf"
	}

	if math.Mod(v, 1.0) == 0 {
		return fmt.Sprintf("%.0f", v)
	}

	return fmt.Sprintf("%f", v)
}

func (o *Number) AsRepr() string {
	return o.AsString()
}

func (o *Number) AsInterface() any {
	return o.Value
}
