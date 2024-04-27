package object

import (
	"fmt"
	"regexp"
	"strings"
)

func _sprintf(args ...Object) string {
	msg := args[0].(*String).Value
	msg = formatInt.ReplaceAllString(msg, `%${1}.0f`)
	v := make([]interface{}, len(args)-1)
	for i, arg := range args[1:] {
		v[i] = arg.AsInterface()
	}
	return fmt.Sprintf(msg, v...)
}

func _sprint(args ...Object) string {
	v := make([]string, len(args))
	for i, arg := range args {
		v[i] = arg.AsString()
	}
	return strings.Join(v, " ")

}

var formatInt, _ = regexp.Compile(`%(-?\d*)d`)
var Printf = F(
	func(scope *Scope, args ...Object) Object {
		fmt.Print(_sprintf(args...))
		return NewTuple(args...)
	},
	`printf`,
	`Prints the formatted string to the standard output.`,
	P("format", V.Type(StringId)),
	P("values").AsSpread(),
)

var Printfln = F(
	func(scope *Scope, args ...Object) Object {
		fmt.Print(_sprintf(args...))
		fmt.Println()
		return NewTuple(args...)
	},
	`printfln`,
	`Prints the formatted string to the standard output.`,
	P("format", V.Type(StringId)),
	P("values").AsSpread(),
)

var Sprintf = F(
	func(scope *Scope, args ...Object) Object {
		return NewString(_sprintf(args...))
	},
	`sprintf`,
	`Returns the formatted string.`,
	P("format", V.Type(StringId)),
	P("values").AsSpread(),
)

var Sprintfln = F(
	func(scope *Scope, args ...Object) Object {
		return NewString(_sprintf(args...) + "\n")
	},
	`sprintfln`,
	`Returns the formatted string.`,
	P("format", V.Type(StringId)),
	P("values").AsSpread(),
)

var Print = F(
	func(scope *Scope, args ...Object) Object {
		fmt.Print(_sprint(args...))
		return NewTuple(args...)
	},
	`print`,
	`Prints the values to the standard output.`,
	P("values").AsSpread(),
)

var Println = F(
	func(scope *Scope, args ...Object) Object {
		fmt.Println(_sprint(args...))
		return NewTuple(args...)
	},
	`println`,
	`Prints the values to the standard output.`,
	P("values").AsSpread(),
)

var Sprint = F(
	func(scope *Scope, args ...Object) Object {
		return NewString(_sprint(args...))
	},
	`sprint`,
	`Returns the formatted string.`,
	P("values").AsSpread(),
)

var Sprintln = F(
	func(scope *Scope, args ...Object) Object {
		return NewString(_sprint(args...))
	},
	`sprintln`,
	`Returns the formatted string.`,
	P("values").AsSpread(),
)

var Range = F(
	func(scope *Scope, args ...Object) Object {
		start := 0.
		end := 0.
		step := 1.

		switch len(args) {
		case 1:
			end = args[0].(*Number).Value
		case 2:
			start = args[0].(*Number).Value
			end = args[1].(*Number).Value
			if start > end {
				step = -1
			}
		case 3:
			start = args[0].(*Number).Value
			end = args[1].(*Number).Value
			step = args[2].(*Number).Value
		}

		cur := start
		inv := step < 0
		return NewInternalStream(func(s *Scope) Object {
			if cur >= end && !inv || cur <= end && inv {
				return nil
			}
			r := cur
			cur += step
			return YieldWith(NewNumber(r))
		}, scope)
	},
	`range`,
	`Returns a range object.`,
	P("values", V.Type(NumberId)).AsSpread(),
)

var Filter = F(
	func(scope *Scope, args ...Object) Object {
		s := StreamTypeObj.Convert(scope, args[0])
		if isRaise(s) {
			return s
		}
		stream := s.(*Stream)

		f := args[1].(*Function)
		return NewInternalStream(func(s *Scope) Object {
			for {
				maybe := Stream_Next.Call(s, stream)
				if isRaise(maybe) {
					return maybe
				}
				if stream.Finished {
					return nil
				}

				value := maybe.(*Maybe).Value
				ret := scope.Eval().Call(scope, f, toLambdaParams(value))
				if isRaise(ret) {
					return ret
				}
				if ret.AsBool() {
					return YieldWith(value)
				}
			}
		}, scope)
	},
	`filter`,
	`Filters the stream.`,
	P("stream"),
	P("f", V.Type(FunctionId)),
)

var Each = F(
	func(scope *Scope, args ...Object) Object {
		s := StreamTypeObj.Convert(scope, args[0])
		if isRaise(s) {
			return s
		}
		stream := s.(*Stream)

		f := args[1].(*Function)
		return NewInternalStream(func(s *Scope) Object {
			maybe := Stream_Next.Call(s, stream)
			if isRaise(maybe) {
				return maybe
			}
			if stream.Finished {
				return nil
			}

			value := maybe.(*Maybe).Value
			ret := scope.Eval().Call(scope, f, toLambdaParams(value))
			if isRaise(ret) {
				return ret
			}
			return YieldWith(value)
		}, scope)
	},
	`each`,
	`Iterates over the stream.`,
	P("stream"),
	P("f", V.Type(FunctionId)),
)

var Map = F(
	func(scope *Scope, args ...Object) Object {
		s := StreamTypeObj.Convert(scope, args[0])
		if isRaise(s) {
			return s
		}
		stream := s.(*Stream)

		f := args[1].(*Function)
		return NewInternalStream(func(s *Scope) Object {
			maybe := Stream_Next.Call(s, stream)
			if isRaise(maybe) {
				return maybe
			}
			if stream.Finished {
				return nil
			}

			value := maybe.(*Maybe).Value
			ret := scope.Eval().Call(scope, f, toLambdaParams(value))
			if isRaise(ret) {
				return ret
			}
			return YieldWith(ret)
		}, scope)
	},
	`map`,
	`Maps the stream.`,
	P("stream"),
	P("f", V.Type(FunctionId)),
)

var Reduce = F(
	func(scope *Scope, args ...Object) Object {
		s := StreamTypeObj.Convert(scope, args[0])
		if isRaise(s) {
			return s
		}
		stream := s.(*Stream)

		f := args[2].(*Function)
		acc := args[1]
		for {
			maybe := Stream_Next.Call(scope, stream)
			if isRaise(maybe) {
				return maybe
			}
			if stream.Finished {
				return acc
			}

			value := maybe.(*Maybe).Value
			acc = scope.Eval().Call(scope, f, toLambdaParams(acc, value))
			if isRaise(acc) {
				return acc
			}
		}
	},
	`reduce`,
	`Reduces the stream.`,
	P("stream"),
	P("acc"),
	P("f", V.Type(FunctionId)),
)

var Sum = F(
	func(scope *Scope, args ...Object) Object {
		s := StreamTypeObj.Convert(scope, args[0])
		if isRaise(s) {
			return s
		}
		stream := s.(*Stream)

		sum := 0.
		for {
			maybe := Stream_Next.Call(scope, stream)
			if isRaise(maybe) {
				return maybe
			}
			if stream.Finished {
				return NewNumber(sum)
			}

			value := maybe.(*Maybe).Value
			if v, ok := value.(*Tuple); ok {
				value = v.Elements[0]
			}

			number, ok := value.(*Number)
			if !ok {
				return scope.Interrupt(Raise("expected number, got %s", value.Type()))
			}
			sum += number.Value
		}
	},
	`sum`,
	`Sums the stream.`,
	P("stream"),
)

var SumBy = F(
	func(scope *Scope, args ...Object) Object {
		s := StreamTypeObj.Convert(scope, args[0])
		if isRaise(s) {
			return s
		}
		stream := s.(*Stream)

		f := args[1].(*Function)
		sum := 0.
		for {
			maybe := Stream_Next.Call(scope, stream)
			if isRaise(maybe) {
				return maybe
			}
			if stream.Finished {
				return NewNumber(sum)
			}

			value := maybe.(*Maybe).Value
			ret := scope.Eval().Call(scope, f, toLambdaParams(value))
			if isRaise(ret) {
				return ret
			}

			number, ok := ret.(*Number)
			if !ok {
				return scope.Interrupt(Raise("expected number, got %s", value.Type()))
			}
			sum += number.Value
		}
	},
	`sumBy`,
	`Sums the stream.`,
	P("stream"),
	P("f", V.Type(FunctionId)),
)

var Count = F(
	func(scope *Scope, args ...Object) Object {
		s := StreamTypeObj.Convert(scope, args[0])
		if isRaise(s) {
			return s
		}
		stream := s.(*Stream)

		count := 0
		for {
			maybe := Stream_Next.Call(scope, stream)
			if isRaise(maybe) {
				return maybe
			}
			if stream.Finished {
				return NewNumber(float64(count))
			}

			count++
		}
	},
	`count`,
	`Counts the stream.`,
	P("stream"),
)

var CountBy = F(
	func(scope *Scope, args ...Object) Object {
		s := StreamTypeObj.Convert(scope, args[0])
		if isRaise(s) {
			return s
		}
		stream := s.(*Stream)

		f := args[1].(*Function)
		count := 0
		for {
			maybe := Stream_Next.Call(scope, stream)
			if isRaise(maybe) {
				return maybe
			}
			if stream.Finished {
				return NewNumber(float64(count))
			}

			value := maybe.(*Maybe).Value
			ret := scope.Eval().Call(scope, f, toLambdaParams(value))
			if isRaise(ret) {
				return ret
			}

			number, ok := ret.(*Number)
			if !ok {
				return scope.Interrupt(Raise("expected number, got %s", value.Type()))
			}
			count += int(number.Value)
		}
	},
	`countBy`,
	`Counts the stream.`,
	P("stream"),
	P("f", V.Type(FunctionId)),
)

func toLambdaParams(values ...Object) []Object {
	params := []Object{}

	for _, value := range values {
		switch value := value.(type) {
		case *Tuple:
			for _, el := range value.Elements {
				params = append(params, el)
			}
		default:
			params = append(params, value)
		}
	}

	return params
}
