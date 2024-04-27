package object

import "math"

func init() {
	NumberTypeObj.AddMethod(Number_Abs)
	NumberTypeObj.AddMethod(Number_Ceil)
	NumberTypeObj.AddMethod(Number_Floor)
	NumberTypeObj.AddMethod(Number_Round)
	NumberTypeObj.AddMethod(Number_RoundToEven)
	NumberTypeObj.AddMethod(Number_Sign)
	NumberTypeObj.AddMethod(Number_CopySign)
	NumberTypeObj.AddMethod(Number_Truncate)
	NumberTypeObj.AddMethod(Number_Int)
	NumberTypeObj.AddMethod(Number_Clamp)
	NumberTypeObj.AddMethod(Number_Remainder)
	NumberTypeObj.AddMethod(Number_Min)
	NumberTypeObj.AddMethod(Number_Max)
	NumberTypeObj.AddMethod(Number_IsOdd)
	NumberTypeObj.AddMethod(Number_IsEven)
	NumberTypeObj.AddMethod(Number_Add)
	NumberTypeObj.AddMethod(Number_Sub)
	NumberTypeObj.AddMethod(Number_Mul)
	NumberTypeObj.AddMethod(Number_Div)
	NumberTypeObj.AddMethod(Number_Mod)
	NumberTypeObj.AddMethod(Number_Pow)
}

var Number_Abs = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		return NewNumber(math.Abs(this.Value))
	},
	`Abs`,
	`Returns the absolute value of a number.`,
	P("this", V.Type(NumberId)),
)

var Number_Ceil = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		return NewNumber(math.Ceil(this.Value))
	},
	`Ceil`,
	`Returns the smallest integer value greater than or equal to a number.`,
	P("this", V.Type(NumberId)),
)

var Number_Floor = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		return NewNumber(math.Floor(this.Value))
	},
	`Floor`,
	`Returns the largest integer value less than or equal to a number.`,
	P("this", V.Type(NumberId)),
)

var Number_Round = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		return NewNumber(math.Round(this.Value))
	},
	`Round`,
	`Returns the nearest integer value to a number.`,
	P("this", V.Type(NumberId)),
)

var Number_RoundToEven = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		return NewNumber(math.RoundToEven(this.Value))
	},
	`RoundToEven`,
	`Returns the nearest even integer value to a number.`,
	P("this", V.Type(NumberId)),
)

var Number_Sign = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		if this.Value >= 0 {
			return One
		}
		return MinusOne
	},
	`Sign`,
	`Returns the sign of a number, 1 for positive, -1 for negative.`,
	P("this", V.Type(NumberId)),
)

var Number_CopySign = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		sign := args[1].(*Number)
		return NewNumber(math.Copysign(this.Value, sign.Value))
	},
	`CopySign`,
	`Returns the value of the first number with the sign of the second number.`,
	P("this", V.Type(NumberId)),
	P("sign", V.Type(NumberId)),
)

var Number_Truncate = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		return NewNumber(math.Trunc(this.Value))
	},
	`Truncate`,
	`Returns the integer value of a number.`,
	P("this", V.Type(NumberId)),
)

var Number_Int = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		return NewNumber(math.Trunc(this.Value))
	},
	`Int`,
	`Returns the integer value of a number.`,
	P("this", V.Type(NumberId)),
)

var Number_Clamp = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		min := args[1].(*Number)
		max := args[2].(*Number)
		return NewNumber(math.Min(math.Max(this.Value, min.Value), max.Value))
	},
	`Clamp`,
	`Clamps a number between a minimum and maximum value.`,
	P("this", V.Type(NumberId)),
	P("min", V.Type(NumberId)),
	P("max", V.Type(NumberId)),
)

var Number_Remainder = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		other := args[1].(*Number)
		return NewNumber(math.Remainder(this.Value, other.Value))
	},
	`Remainder`,
	`Returns the remainder of a division between two numbers.`,
	P("this", V.Type(NumberId)),
	P("other", V.Type(NumberId)),
)

var Number_Min = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		for _, arg := range args[1:] {
			other := arg.(*Number)
			if other.Value < this.Value {
				this = other
			}
		}
		return this
	},
	`Min`,
	`Returns the smallest number from a list of numbers.`,
	P("this", V.Type(NumberId)),
	P("others", V.Type(NumberId)).AsSpread(),
)

var Number_Max = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		for _, arg := range args[1:] {
			other := arg.(*Number)
			if other.Value > this.Value {
				this = other
			}
		}
		return this
	},
	`Max`,
	`Returns the largest number from a list of numbers.`,
	P("this", V.Type(NumberId)),
	P("others", V.Type(NumberId)).AsSpread(),
)

var Number_IsOdd = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		return NewBoolean(int64(this.Value)%2 != 0)
	},
	`IsOdd`,
	`Returns true if the number is odd.`,
	P("this", V.Type(NumberId)),
)

var Number_IsEven = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		return NewBoolean(int64(this.Value)%2 == 0)
	},
	`IsEven`,
	`Returns true if the number is even.`,
	P("this", V.Type(NumberId)),
)

var Number_Add = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		other := args[1].(*Number)
		return NewNumber(this.Value + other.Value)
	},
	`Add`,
	`Returns the sum of two numbers.`,
	P("this", V.Type(NumberId)),
	P("other", V.Type(NumberId)),
)

var Number_Sub = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		other := args[1].(*Number)
		return NewNumber(this.Value - other.Value)
	},
	`Sub`,
	`Returns the difference of two numbers.`,
	P("this", V.Type(NumberId)),
	P("other", V.Type(NumberId)),
)

var Number_Mul = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		other := args[1].(*Number)
		return NewNumber(this.Value * other.Value)
	},
	`Mul`,
	`Returns the product of two numbers.`,
	P("this", V.Type(NumberId)),
	P("other", V.Type(NumberId)),
)

var Number_Div = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		other := args[1].(*Number)
		return NewNumber(this.Value / other.Value)
	},
	`Div`,
	`Returns the division of two numbers.`,
	P("this", V.Type(NumberId)),
	P("other", V.Type(NumberId)),
)

var Number_Mod = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		other := args[1].(*Number)
		return NewNumber(math.Mod(this.Value, other.Value))
	},
	`Mod`,
	`Returns the remainder of a division between two numbers.`,
	P("this", V.Type(NumberId)),
	P("other", V.Type(NumberId)),
)

var Number_Pow = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Number)
		other := args[1].(*Number)
		return NewNumber(math.Pow(this.Value, other.Value))
	},
	`Pow`,
	`Returns the first number raised to the power of the second number.`,
	P("this", V.Type(NumberId)),
	P("other", V.Type(NumberId)),
)
