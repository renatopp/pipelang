package object

import (
	"sort"
	"strings"
)

func init() {
	ListTypeObj.AddMethod(List_Set)
	ListTypeObj.AddMethod(List_Push)
	ListTypeObj.AddMethod(List_Pop)
	ListTypeObj.AddMethod(List_Insert)
	ListTypeObj.AddMethod(List_Remove)
	ListTypeObj.AddMethod(List_RemoveAt)
	ListTypeObj.AddMethod(List_Copy)
	ListTypeObj.AddMethod(List_Clear)
	ListTypeObj.AddMethod(List_Concat)
	ListTypeObj.AddMethod(List_Split)
	ListTypeObj.AddMethod(List_SplitAt)
	ListTypeObj.AddMethod(List_SplitFn)
	ListTypeObj.AddMethod(List_Size)
	ListTypeObj.AddMethod(List_Sorted)
	ListTypeObj.AddMethod(List_Sort)
	ListTypeObj.AddMethod(List_SortedFn)
	ListTypeObj.AddMethod(List_SortFn)
	ListTypeObj.AddMethod(List_Reversed)
	ListTypeObj.AddMethod(List_Reverse)
	ListTypeObj.AddMethod(List_Get)
	ListTypeObj.AddMethod(List_GetOr)
	ListTypeObj.AddMethod(List_Sub)
	ListTypeObj.AddMethod(List_Find)
	ListTypeObj.AddMethod(List_FindFn)
	ListTypeObj.AddMethod(List_FindLast)
	ListTypeObj.AddMethod(List_FindLastFn)
	ListTypeObj.AddMethod(List_FindAll)
	ListTypeObj.AddMethod(List_FindAllFn)
	ListTypeObj.AddMethod(List_Contains)
	ListTypeObj.AddMethod(List_ContainsFn)
	ListTypeObj.AddMethod(List_IsEmpty)
	ListTypeObj.AddMethod(List_Count)
	ListTypeObj.AddMethod(List_CountFn)
	ListTypeObj.AddMethod(List_Join)
	ListTypeObj.AddMethod(List_Elements)
}

// ----------------------------------------------------------------------------
// Add/Remove
// ----------------------------------------------------------------------------
var List_Set = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		index := int(args[1].(*Number).Value)
		len := len(this.Elements)
		if index < 0 {
			index = len + index
		}

		if index < 0 || index >= len {
			return scope.Interrupt(Raise("index out of range"))
		}

		this.Elements[index] = args[2]
		return this
	},
	`Set`,
	`Sets the element at the specified index.`,
	P("this", V.Type(ListId)),
	P("index", V.Type(NumberId)),
	P("element"),
)

var List_Push = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		this.Elements = append(this.Elements, args[1])
		return this
	},
	`Push`,
	`Appends an element to the end of the list.`,
	P("this", V.Type(ListId)),
	P("element"),
)

var List_Pop = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		if len(this.Elements) == 0 {
			return scope.Interrupt(Raise("cannot pop from an empty list"))
		}

		last := this.Elements[len(this.Elements)-1]
		this.Elements = this.Elements[:len(this.Elements)-1]
		return last
	},
	`Pop`,
	`Removes the last element from the list and returns it.`,
	P("this", V.Type(ListId)),
)

var List_Insert = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		index := int(args[1].(*Number).Value)
		element := args[2]
		len := len(this.Elements)

		if index < 0 {
			index = len + index
		}

		if index < 0 || index > len {
			return scope.Interrupt(Raise("index out of range"))
		}

		this.Elements = append(this.Elements[:index], append([]Object{element}, this.Elements[index:]...)...)
		return this
	},
	`Insert`,
	`Inserts an element at the specified index.`,
	P("this", V.Type(ListId)),
	P("index", V.Type(NumberId)),
	P("element"),
)

var List_Remove = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		element := args[1]
		for i, e := range this.Elements {
			if scope.Eval().Operator(scope, "==", e, element) == True {
				this.Elements = append(this.Elements[:i], this.Elements[i+1:]...)
				return this
			}
		}
		return scope.Interrupt(Raise("element not found"))
	},
	`Remove`,
	`Removes the first occurrence of the specified element from the list.`,
	P("this", V.Type(ListId)),
	P("element"),
)

var List_RemoveAt = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		index := int(args[1].(*Number).Value)
		len := len(this.Elements)

		if index < 0 {
			index = len + index
		}

		if index < 0 || index >= len {
			return scope.Interrupt(Raise("index out of range"))
		}

		this.Elements = append(this.Elements[:index], this.Elements[index+1:]...)
		return this
	},
	`RemoveAt`,
	`Removes the element at the specified index.`,
	P("this", V.Type(ListId)),
	P("index", V.Type(NumberId)),
)

var List_Copy = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		elements := make([]Object, len(this.Elements))
		copy(elements, this.Elements)
		return NewList(elements...)
	},
	`Copy`,
	`Returns a shallow copy of the list.`,
	P("this", V.Type(ListId)),
)

var List_Clear = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		this.Elements = this.Elements[:0]
		return this
	},
	`Clear`,
	`Removes all elements from the list.`,
	P("this", V.Type(ListId)),
)

var List_Concat = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		others := args[1:]
		e := this.Elements
		for _, other := range others {
			e = append(e, other.(*List).Elements...)
		}
		return NewList(e...)
	},
	`Concat`,
	`Concatenates the list with one or more other lists.`,
	P("this", V.Type(ListId)),
	P("others", V.Type(ListId)).AsSpread(),
)

var List_Split = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		separator := args[1]
		result := NewList()
		sublist := NewList()
		for _, e := range this.Elements {
			if scope.Eval().Operator(scope, "==", e, separator) == True {
				result.Elements = append(result.Elements, sublist)
				sublist = NewList()
			} else {
				sublist.Elements = append(sublist.Elements, e)
			}
		}
		result.Elements = append(result.Elements, sublist)
		return result
	},
	`Split`,
	`Splits the list into sublists based on the separator.`,
	P("this", V.Type(ListId)),
	P("separator"),
)

var List_SplitAt = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		index := int(args[1].(*Number).Value)
		len := len(this.Elements)

		if index < 0 {
			index = len + index
		}

		if index < 0 || index >= len {
			return scope.Interrupt(Raise("index out of range"))
		}
		return NewList(NewList(this.Elements[:index]...), NewList(this.Elements[index:]...))
	},
	`SplitAt`,
	`Splits the list into two sublists at the specified index.`,
	P("this", V.Type(ListId)),
	P("index", V.Type(NumberId)),
)

var List_SplitFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		f := args[1].(*Function)

		result := NewList()
		sublist := NewList()
		for i, e := range this.Elements {
			ret := scope.Eval().Call(scope, f, []Object{e, NewNumber(float64(i))})
			if isRaise(ret) {
				return ret
			}

			if ret.AsBool() {
				result.Elements = append(result.Elements, sublist)
				sublist = NewList()
			}
			sublist.Elements = append(sublist.Elements, e)
		}
		if len(sublist.Elements) > 0 {
			result.Elements = append(result.Elements, sublist)
		}
		return result
	},
	`SplitFn`,
	`Splits the list into sublists based on the result of the function.`,
	P("this", V.Type(ListId)),
	P("f", V.Type(FunctionId)),
)

// ----------------------------------------------------------------------------
// Ordering
// ----------------------------------------------------------------------------
var List_Sorted = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		sort.Slice(this.Elements, func(i, j int) bool {
			ret := scope.Eval().Operator(scope, "<", this.Elements[i], this.Elements[j])
			if isRaise(ret) {
				return false
			}
			return ret.AsBool()
		})
		return this
	},
	`Sorted`,
	`Sorts the list in ascending order.`,
	P("this", V.Type(ListId)),
)

var List_Sort = F(
	func(scope *Scope, args ...Object) Object {
		copy := List_Copy.Call(scope, args[0])
		if isRaise(copy) {
			return copy
		}
		return List_Sorted.Call(scope, copy)
	},
	`Sort`,
	`Copy and sorts the new list in ascending order.`,
	P("this", V.Type(ListId)),
)

var List_SortedFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		f := args[1].(*Function)
		sort.Slice(this.Elements, func(i, j int) bool {
			ret := scope.Eval().Call(scope, f, []Object{this.Elements[i], this.Elements[j]})
			if isRaise(ret) {
				return false
			}
			return ret.AsBool()
		})
		return this
	},
	`SortedFn`,
	`Sorts the list based on the result of the function.`,
	P("this", V.Type(ListId)),
	P("f", V.Type(FunctionId)),
)

var List_SortFn = F(
	func(scope *Scope, args ...Object) Object {
		copy := List_Copy.Call(scope, args[0])
		if isRaise(copy) {
			return copy
		}
		return List_SortedFn.Call(scope, copy, args[1])
	},
	`SortFn`,
	`Copy and sorts the new list based on the result of the function.`,
	P("this", V.Type(ListId)),
	P("f", V.Type(FunctionId)),
)

var List_Reversed = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		for i, j := 0, len(this.Elements)-1; i < j; i, j = i+1, j-1 {
			this.Elements[i], this.Elements[j] = this.Elements[j], this.Elements[i]
		}
		return this
	},
	`Reversed`,
	`Reverses the order of the elements in the list.`,
	P("this", V.Type(ListId)),
)

var List_Reverse = F(
	func(scope *Scope, args ...Object) Object {
		copy := List_Copy.Call(scope, args[0])
		if isRaise(copy) {
			return copy
		}
		return List_Reversed.Call(scope, copy)
	},
	`Reverse`,
	`Copy and reverses the order of the elements in the new list.`,
	P("this", V.Type(ListId)),
)

// ----------------------------------------------------------------------------
// Indexing and Search
// ----------------------------------------------------------------------------
var List_Get = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		index := int(args[1].(*Number).Value)
		len := len(this.Elements)

		if index < 0 {
			index = len + index
		}

		if index < 0 || index >= len {
			return scope.Interrupt(Raise("index out of range"))
		}

		return this.Elements[index]
	},
	`Get`,
	`Returns the element at the specified index.`,
	P("this", V.Type(ListId)),
	P("index", V.Type(NumberId)),
)

var List_GetOr = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		index := int(args[1].(*Number).Value)
		len := len(this.Elements)

		if index < 0 {
			index = len + index
		}

		if index < 0 || index >= len {
			return args[2]
		}

		return this.Elements[index]
	},
	`GetOr`,
	`Returns the element at the specified index or the default value if the index is out of range.`,
	P("this", V.Type(ListId)),
	P("index", V.Type(NumberId)),
	P("default"),
)

var List_Sub = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		from := int(args[1].(*Number).Value)
		to := int(args[2].(*Number).Value)
		len := len(this.Elements)

		if from < 0 {
			from = len + from
		}

		if to < 0 {
			to = len + to
		}

		if to < from || to < 0 || from >= len {
			return NewList()
		}

		from = max(0, from)
		to = min(len, to)
		return NewList(this.Elements[from:to]...)
	},
	`Sub`,
	`Returns a new list containing the elements from the start index to the end index.`,
	P("this", V.Type(ListId)),
	P("start", V.Type(NumberId)),
	P("end", V.Type(NumberId)),
)

var List_Find = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		targets := args[1:]
		for i, e := range this.Elements {
			for _, target := range targets {
				ret := scope.Eval().Operator(scope, "==", e, target)
				if isRaise(ret) {
					return ret
				} else if ret.AsBool() {
					return NewNumber(float64(i))
				}
			}
		}
		return MinusOne
	},
	`Find`,
	`Returns the index of the first occurrence of any of the specified elements in the list.`,
	P("this", V.Type(ListId)),
	P("targets").AsSpread(),
)

var List_FindFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		f := args[1].(*Function)
		for i, e := range this.Elements {
			ret := scope.Eval().Call(scope, f, []Object{e, NewNumber(float64(i))})
			if isRaise(ret) {
				return ret
			} else if ret.AsBool() {
				return NewNumber(float64(i))
			}
		}
		return MinusOne
	},
	`FindFn`,
	`Returns the index of the first element that satisfies the function.`,
	P("this", V.Type(ListId)),
	P("f", V.Type(FunctionId)),
)

var List_FindLast = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		targets := args[1:]
		for i := len(this.Elements) - 1; i >= 0; i-- {
			e := this.Elements[i]
			for _, target := range targets {
				ret := scope.Eval().Operator(scope, "==", e, target)
				if isRaise(ret) {
					return ret
				} else if ret.AsBool() {
					return NewNumber(float64(i))
				}
			}
		}
		return MinusOne
	},
	`FindLast`,
	`Returns the index of the last occurrence of any of the specified elements in the list.`,
	P("this", V.Type(ListId)),
	P("targets").AsSpread(),
)

var List_FindLastFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		f := args[1].(*Function)
		for i := len(this.Elements) - 1; i >= 0; i-- {
			e := this.Elements[i]
			ret := scope.Eval().Call(scope, f, []Object{e, NewNumber(float64(i))})
			if isRaise(ret) {
				return ret
			} else if ret.AsBool() {
				return NewNumber(float64(i))
			}
		}
		return MinusOne
	},
	`FindLastFn`,
	`Returns the index of the last element that satisfies the function.`,
	P("this", V.Type(ListId)),
	P("f", V.Type(FunctionId)),
)

var List_FindAll = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		targets := args[1:]
		result := NewList()
		for i, e := range this.Elements {
			for _, target := range targets {
				ret := scope.Eval().Operator(scope, "==", e, target)
				if isRaise(ret) {
					return ret
				} else if ret.AsBool() {
					result.Elements = append(result.Elements, NewNumber(float64(i)))
					break
				}
			}
		}
		return result
	},
	`FindAll`,
	`Returns the indices of all occurrences of any of the specified elements in the list.`,
	P("this", V.Type(ListId)),
	P("targets").AsSpread(),
)

var List_FindAllFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		f := args[1].(*Function)
		result := NewList()
		for i, e := range this.Elements {
			ret := scope.Eval().Call(scope, f, []Object{e, NewNumber(float64(i))})
			if isRaise(ret) {
				return ret
			} else if ret.AsBool() {
				result.Elements = append(result.Elements, NewNumber(float64(i)))
			}
		}
		return result
	},
	`FindAllFn`,
	`Returns the indices of all elements that satisfy the function.`,
	P("this", V.Type(ListId)),
	P("f", V.Type(FunctionId)),
)

// ----------------------------------------------------------------------------
// Checkers
// ----------------------------------------------------------------------------
var List_Size = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		return NewNumber(float64(len(this.Elements)))
	},
	`Size`,
	`Returns the number of elements in the list.`,
	P("this", V.Type(ListId)),
)

var List_Contains = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		target := args[1]
		for _, e := range this.Elements {
			ret := scope.Eval().Operator(scope, "==", e, target)
			if isRaise(ret) {
				return ret
			} else if ret.AsBool() {
				return True
			}
		}
		return False
	},
	`Contains`,
	`Returns true if the list contains the specified element.`,
	P("this", V.Type(ListId)),
	P("element"),
)

var List_ContainsFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		f := args[1].(*Function)
		for _, e := range this.Elements {
			ret := scope.Eval().Call(scope, f, []Object{e})
			if isRaise(ret) {
				return ret
			} else if ret.AsBool() {
				return True
			}
		}
		return False
	},
	`ContainsFn`,
	`Returns true if the list contains an element that satisfies the function.`,
	P("this", V.Type(ListId)),
	P("f", V.Type(FunctionId)),
)

var List_IsEmpty = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		return NewBoolean(len(this.Elements) == 0)
	},
	`IsEmpty`,
	`Returns true if the list is empty.`,
	P("this", V.Type(ListId)),
)

var List_Count = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		targets := args[1:]
		count := 0
		for _, e := range this.Elements {
			for _, target := range targets {
				ret := scope.Eval().Operator(scope, "==", e, target)
				if isRaise(ret) {
					return ret
				} else if ret.AsBool() {
					count++
					break
				}
			}
		}
		return NewNumber(float64(count))
	},
	`Count`,
	`Returns the number of occurrences of any of the specified elements in the list.`,
	P("this", V.Type(ListId)),
	P("targets").AsSpread(),
)

var List_CountFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		f := args[1].(*Function)
		count := 0
		for _, e := range this.Elements {
			ret := scope.Eval().Call(scope, f, []Object{e})
			if isRaise(ret) {
				return ret
			} else if ret.AsBool() {
				count++
			}
		}
		return NewNumber(float64(count))
	},
	`CountFn`,
	`Returns the number of elements that satisfy the function.`,
	P("this", V.Type(ListId)),
	P("f", V.Type(FunctionId)),
)

// ----------------------------------------------------------------------------
// Formatter
// ----------------------------------------------------------------------------
var List_Join = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		sep := args[1].(*String).Value
		var elements []string
		for _, e := range this.Elements {
			elements = append(elements, e.AsString())
		}
		return NewString(strings.Join(elements, sep))
	},
	`Join`,
	`Concatenates the elements of the list into a single string using the specified separator.`,
	P("this", V.Type(ListId)),
	P("separator", V.Type(StringId)),
)

// ----------------------------------------------------------------------------
// Stream
// ----------------------------------------------------------------------------
var List_Elements = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*List)
		idx := 0
		return NewInternalStream(func(s *Scope) Object {
			if idx >= len(this.Elements) {
				return nil
			}
			e := this.Elements[idx]
			idx++
			return YieldWith(NewTuple(e, NewNumber(float64(idx-1))))
		}, scope)
	},
	`Elements`,
	`Returns a stream of the elements in the list.`,
	P("this", V.Type(ListId)),
)
