package object

func init() {
	DictTypeObj.AddMethod(Dict_Set)
	DictTypeObj.AddMethod(Dict_Get)
	DictTypeObj.AddMethod(Dict_GetOr)
	DictTypeObj.AddMethod(Dict_Has)
	DictTypeObj.AddMethod(Dict_Remove)
	DictTypeObj.AddMethod(Dict_Contains)
	DictTypeObj.AddMethod(Dict_ContainsFn)
	DictTypeObj.AddMethod(Dict_Size)
	DictTypeObj.AddMethod(Dict_Clear)
	DictTypeObj.AddMethod(Dict_Copy)
	DictTypeObj.AddMethod(Dict_Concat)
	DictTypeObj.AddMethod(Dict_Find)
	DictTypeObj.AddMethod(Dict_FindFn)
	DictTypeObj.AddMethod(Dict_FindAll)
	DictTypeObj.AddMethod(Dict_FindAllFn)
	DictTypeObj.AddMethod(Dict_IsEmpty)
	DictTypeObj.AddMethod(Dict_Count)
	DictTypeObj.AddMethod(Dict_CountFn)
	DictTypeObj.AddMethod(Dict_Keys)
	DictTypeObj.AddMethod(Dict_Values)
	DictTypeObj.AddMethod(Dict_Items)
	DictTypeObj.AddMethod(Dict_Elements)
}

var Dict_Set = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		index := args[1]
		this.Elements[index.AsString()] = args[2]
		return this
	},
	`Set`,
	`Sets the element at the specified index.`,
	P("this", V.Type(DictId)),
	P("index"),
	P("element"),
)

var Dict_Get = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		index := args[1]
		res, ok := this.Elements[index.AsString()]
		if !ok {
			return scope.Interrupt(Raise("key not found: %s", index.AsString()))
		}
		return res
	},
	`Get`,
	`Gets the element at the specified index.`,
	P("this", V.Type(DictId)),
	P("index"),
)

var Dict_GetOr = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		index := args[1]
		def := args[2]
		res, ok := this.Elements[index.AsString()]
		if !ok {
			return def
		}
		return res
	},
	`GetOr`,
	`Gets the element at the specified index or returns the default value.`,
	P("this", V.Type(DictId)),
	P("index"),
	P("default"),
)

var Dict_Has = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		index := args[1]
		_, ok := this.Elements[index.AsString()]
		return NewBoolean(ok)
	},
	`Has`,
	`Returns true if the specified index exists in the dictionary.`,
	P("this", V.Type(DictId)),
	P("index"),
)

var Dict_Remove = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		index := args[1]
		res, ok := this.Elements[index.AsString()]
		if !ok {
			return scope.Interrupt(Raise("key not found: %s", index.AsString()))
		}
		delete(this.Elements, index.AsString())
		return res
	},
	`Remove`,
	`Removes the element at the specified index.`,
	P("this", V.Type(DictId)),
	P("index"),
)

var Dict_Contains = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		element := args[1]
		for _, v := range this.Elements {
			if scope.eval.Operator(scope, "==", v, element).AsBool() {
				return True
			}
		}
		return False
	},
	`Contains`,
	`Returns true if the specified element exists in the dictionary.`,
	P("this", V.Type(DictId)),
	P("element"),
)

var Dict_ContainsFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		f := args[1]
		for k, v := range this.Elements {
			if scope.eval.Call(scope, f, []Object{v, NewString(k)}).AsBool() {
				return True
			}
		}
		return False
	},
	`ContainsFn`,
	`Returns true if the specified function returns true for any element in the dictionary.`,
	P("this", V.Type(DictId)),
	P("f", V.Type(FunctionId)),
)

var Dict_Size = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		return NewNumber(float64(len(this.Elements)))
	},
	`Size`,
	`Returns the number of elements in the dictionary.`,
	P("this", V.Type(DictId)),
)

var Dict_Clear = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		this.Elements = make(map[string]Object)
		return this
	},
	`Clear`,
	`Removes all elements from the dictionary.`,
	P("this", V.Type(DictId)),
)

var Dict_Copy = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		elements := make(map[string]Object)
		for k, v := range this.Elements {
			elements[k] = v
		}
		return NewDict(elements)
	},
	`Copy`,
	`Returns a copy of the dictionary.`,
	P("this", V.Type(DictId)),
)

var Dict_Concat = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		others := args[1:]
		elements := this.Elements
		for _, other := range others {
			for k, v := range other.(*Dict).Elements {
				elements[k] = v
			}
		}
		return NewDict(elements)
	},
	`Concat`,
	`Returns a new dictionary with the elements of both dictionaries.`,
	P("this", V.Type(DictId)),
	P("other", V.Type(DictId)).AsSpread(),
)

var Dict_Find = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		elements := args[1:]
		for k, v := range this.Elements {
			for _, element := range elements {
				if scope.eval.Operator(scope, "==", v, element).AsBool() {
					return NewString(k)
				}
			}
		}
		return False
	},
	`Find`,
	`Returns the key of the first element that matches the specified element.`,
	P("this", V.Type(DictId)),
	P("element").AsSpread(),
)

var Dict_FindFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		f := args[1]
		for k, v := range this.Elements {
			ko := NewString(k)
			if scope.eval.Call(scope, f, []Object{v, ko}).AsBool() {
				return ko
			}
		}
		return False
	},
	`FindFn`,
	`Returns the key of the first element that matches the specified function.`,
	P("this", V.Type(DictId)),
	P("f", V.Type(FunctionId)),
)

var Dict_FindAll = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		elements := args[1:]
		var keys []Object
		for k, v := range this.Elements {
			for _, element := range elements {
				if scope.eval.Operator(scope, "==", v, element).AsBool() {
					keys = append(keys, NewString(k))
				}
			}
		}
		return NewList(keys...)
	},
	`FindAll`,
	`Returns the keys of all elements that match the specified element.`,
	P("this", V.Type(DictId)),
	P("element").AsSpread(),
)

var Dict_FindAllFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		f := args[1]
		var keys []Object
		for k, v := range this.Elements {
			ko := NewString(k)
			if scope.eval.Call(scope, f, []Object{v, ko}).AsBool() {
				keys = append(keys, ko)
			}
		}
		return NewList(keys...)
	},
	`FindAllFn`,
	`Returns the keys of all elements that match the specified function.`,
	P("this", V.Type(DictId)),
	P("f", V.Type(FunctionId)),
)

var Dict_IsEmpty = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		return NewBoolean(len(this.Elements) == 0)
	},
	`IsEmpty`,
	`Returns true if the dictionary is empty.`,
)

var Dict_Count = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		elements := args[1:]
		var count int
		for _, v := range this.Elements {
			for _, element := range elements {
				if scope.eval.Operator(scope, "==", v, element).AsBool() {
					count++
				}
			}
		}
		return NewNumber(float64(count))
	},
	`Count`,
	`Returns the number of elements that match the specified element.`,
	P("this", V.Type(DictId)),
	P("element").AsSpread(),
)

var Dict_CountFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		f := args[1]
		var count int
		for k, v := range this.Elements {
			if scope.eval.Call(scope, f, []Object{v, NewString(k)}).AsBool() {
				count++
			}
		}
		return NewNumber(float64(count))
	},
	`CountFn`,
	`Returns the number of elements that match the specified function.`,
	P("this", V.Type(DictId)),
	P("f", V.Type(FunctionId)),
)

var Dict_Keys = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		var keys []Object
		for k := range this.Elements {
			keys = append(keys, NewString(k))
		}
		return NewList(keys...)
	},
	`Keys`,
	`Returns a list of all keys in the dictionary.`,
	P("this", V.Type(DictId)),
)

var Dict_Values = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		var values []Object
		for _, v := range this.Elements {
			values = append(values, v)
		}
		return NewList(values...)
	},
	`Values`,
	`Returns a list of all values in the dictionary.`,
	P("this", V.Type(DictId)),
)

var Dict_Items = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		var items []Object
		for k, v := range this.Elements {
			items = append(items, NewList(NewString(k), v))
		}
		return NewList(items...)
	},
	`Items`,
	`Returns a list of all key-value pairs in the dictionary.`,
	P("this", V.Type(DictId)),
)

var Dict_Elements = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*Dict)
		keys := []string{}
		for k := range this.Elements {
			keys = append(keys, k)
		}

		idx := 0
		return NewInternalStream(func(s *Scope) Object {
			if idx >= len(keys) {
				return nil
			}
			k := keys[idx]
			idx++
			return YieldWith(NewTuple(this.Elements[k], NewString(k)))
		}, scope)
	},
	`Elements`,
	`Returns a stream of all key-value pairs in the dictionary.`,
	P("this", V.Type(DictId)),
)
