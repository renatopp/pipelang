package evaluator

import (
	i "github.com/renatopp/pipelang/internal"
	"github.com/renatopp/pipelang/internal/ast"
	o "github.com/renatopp/pipelang/internal/object"
)

const ForReturnKey = "$for-return"
const ForInKey = "$for-in"

type Evaluator struct {
	Scope *o.Scope
}

func New(scope *o.Scope) *Evaluator {
	r := &Evaluator{}
	r.Scope = scope
	r.Scope.WithEval(r)

	return r
}

func (r *Evaluator) Eval(node ast.Node) (o.Object, *i.Error) {
	return r.EvalWithScope(r.Scope, node)
}

func (r *Evaluator) EvalWithScope(scope *o.Scope, node ast.Node) (o.Object, *i.Error) {
	value := r.eval(scope, node)

	if intr := asInterruption(value); intr != nil {
		if intr.Category == o.RaiseId {
			return nil, i.NewError(intr.Value.AsString(), intr.Stack)
		}
		value = intr.Value
	}

	return value, nil
}

func (r *Evaluator) RawEval(scope *o.Scope, node ast.Node) o.Object {
	return r.eval(scope, node)
}

func (r *Evaluator) Call(scope *o.Scope, target o.Object, args []o.Object) o.Object {
	return r.call(scope, target, args)
}

func (r *Evaluator) Operator(scope *o.Scope, op string, left, right o.Object) o.Object {
	return r.evalOperator(scope, op, left, right)
}

func (r *Evaluator) eval(scope *o.Scope, node ast.Node) o.Object {
	scope.PushNode(node)
	defer scope.PopNode()

	switch n := node.(type) {

	// Types
	case *ast.Number:
		return r.evalNumber(scope, n)

	case *ast.Boolean:
		return r.evalBoolean(scope, n)

	case *ast.String:
		return r.evalString(scope, n)

	case *ast.Identifier:
		return r.evalIdentifier(scope, n)

	case *ast.Tuple:
		return r.evalTuple(scope, n)

	case *ast.List:
		return r.evalList(scope, n)

	case *ast.Dict:
		return r.evalDict(scope, n)

	case *ast.FunctionDef:
		return r.evalFunctionDef(scope, n)

	case *ast.DataDef:
		return r.evalDataDef(scope, n)

	// Operations
	case *ast.PrefixOperator:
		return r.evalPrefixOperator(scope, n)

	case *ast.InfixOperator:
		return r.evalInfixOperator(scope, n)

	case *ast.Call:
		return r.evalCall(scope, n)

	case *ast.Instantiate:
		return r.evalInstantiate(scope, n)

	case *ast.Assignment:
		return r.evalAssignment(scope, n)

	case *ast.Spread:
		return r.evalSpread(scope, n)

	case *ast.Access:
		return r.evalAccess(scope, n)

	case *ast.Index:
		return r.evalIndex(scope, n)

	case *ast.Wrap:
		return r.evalWrap(scope, n)

	case *ast.Unwrap:
		return r.evalUnwrap(scope, n)

	// Control Flow
	case *ast.Block:
		return r.evalBlock(scope, n)

	case *ast.Return:
		return r.evalReturn(scope, n)

	case *ast.Raise:
		return r.evalRaise(scope, n)

	case *ast.Yield:
		return r.evalYield(scope, n)

	case *ast.Break:
		return r.evalBreak(scope, n)

	case *ast.Continue:
		return r.evalContinue(scope, n)

	case *ast.If:
		return r.evalIf(scope, n)

	case *ast.For:
		return r.evalFor(scope, n)

	case *ast.With:
		return r.evalWith(scope, n)

	case *ast.Match:
		return r.evalMatch(scope, n)

	default:
		return scope.Interrupt(o.Raise("unknown node type '%v'", node))
	}
}

// ----------------------------------------------------------------------------
// Types Evaluation
// ----------------------------------------------------------------------------
func (r *Evaluator) evalNumber(_ *o.Scope, n *ast.Number) o.Object {
	return o.NewNumber(n.Value)
}

func (r *Evaluator) evalBoolean(_ *o.Scope, n *ast.Boolean) o.Object {
	return o.NewBoolean(n.Value)
}

func (r *Evaluator) evalString(_ *o.Scope, n *ast.String) o.Object {
	return o.NewString(n.Value)
}

func (r *Evaluator) evalIdentifier(scope *o.Scope, n *ast.Identifier) o.Object {
	obj := scope.GetGlobal(n.Value)
	if obj == nil {
		return scope.Interrupt(o.Raise("identifier '%s' not found", n.Value))
	}
	obj.SetParent(nil)
	return obj
}

func (r *Evaluator) evalTuple(scope *o.Scope, n *ast.Tuple) o.Object {
	elements := []o.Object{}
	for _, element := range n.Elements {
		_, isSpread := element.(*ast.Spread)

		item := r.eval(scope, element)
		if isRaise(item) {
			return item
		}

		if isSpread { // Deconstruct the spread into the tuple (a...) => (a1, a2, a3)
			tuple := o.TupleTypeObj.Convert(scope, item)
			if isRaise(tuple) {
				return tuple
			}
			elements = append(elements, tuple.(*o.Tuple).Elements...)
		} else {
			elements = append(elements, item)
		}
	}
	return o.NewTuple(elements...)
}

func (r *Evaluator) evalList(scope *o.Scope, n *ast.List) o.Object {
	elements := []o.Object{}
	for _, element := range n.Elements {
		_, isSpread := element.(*ast.Spread)

		item := r.eval(scope, element)
		if isRaise(item) {
			return item
		}

		if isSpread { // Deconstruct the spread into the tuple (a...) => (a1, a2, a3)
			tuple := o.TupleTypeObj.Convert(scope, item)
			if isRaise(tuple) {
				return tuple
			}
			elements = append(elements, tuple.(*o.Tuple).Elements...)
		} else {
			elements = append(elements, item)
		}
	}
	return o.NewList(elements...)
}

func (r *Evaluator) evalDict(scope *o.Scope, n *ast.Dict) o.Object {
	elements := []o.Object{}
	for _, element := range n.Elements {
		item := r.eval(scope, element)
		if isRaise(item) {
			return item
		}

		elements = append(elements, item)
	}
	return o.NewDictFromList(elements...)
}

func (r *Evaluator) evalFunctionDef(scope *o.Scope, n *ast.FunctionDef) o.Object {
	fn := o.NewFunction(n.Name, n.Parameters, n.Body, scope)

	var ret o.Object = fn
	if n.Generator {
		ret = o.NewBuiltinFunction(n.Name, func(scope *o.Scope, args ...o.Object) o.Object {
			genScope := fn.Scope.New()
			params := &ast.Tuple{Elements: fn.Parameters}
			ret := r.resolveAssignment(genScope, ":=", params, o.NewTuple(args...))
			if isRaise(ret) {
				return ret
			}

			return o.NewStream(fn, genScope)
		})
	}

	if n.Name != "" {
		scope.SetLocal(n.Name, ret)
	}
	return ret
}

func (r *Evaluator) evalDataDef(scope *o.Scope, n *ast.DataDef) o.Object {
	attributes := map[string]ast.Node{}
	methods := map[string]*o.Function{}

	for _, ext := range n.Extensions {
		ext := r.eval(scope, ext)
		if isRaise(ext) {
			return ext
		}

		if ext.TypeId() != o.DataId {
			return scope.Interrupt(o.Raise("type '%s' cannot be extended", ext.TypeId()))
		}

		data := ext.(*o.DataType)
		for name, node := range data.Attributes {
			attributes[name] = node
		}
		for name, method := range data.Methods {
			methods[name] = method
		}
	}

	for name, node := range n.Attributes {
		attributes[name] = node
	}
	for name, method := range n.Methods {
		fn := r.eval(scope, method)
		if isRaise(fn) {
			return fn
		}
		methods[name] = fn.(*o.Function)
	}

	data := o.NewDataType(n.Name, attributes, methods)
	if n.Name != "" {
		scope.SetLocal(n.Name, data)
	}

	return data
}

// ----------------------------------------------------------------------------
// Operation Evaluation
// ----------------------------------------------------------------------------
func (r *Evaluator) evalPrefixOperator(scope *o.Scope, n *ast.PrefixOperator) o.Object {
	right := r.eval(scope, n.Right)

	switch n.Operator {
	case "+", "-":
		if right.TypeId() != o.NumberId {
			return scope.Interrupt(o.Raise("type '%s' does not support unary operator '%s'", right.TypeId(), n.Operator))
		}

		right := right.(*o.Number)
		if n.Operator == "+" {
			return o.NewNumber(right.Value)
		}
		return o.NewNumber(-right.Value)

	case "not":
		return o.NewBoolean(!right.AsBool())

	default:
		return scope.Interrupt(o.Raise("unknown unary operator '%s'", n.Operator))
	}
}

func (r *Evaluator) evalInfixOperator(scope *o.Scope, n *ast.InfixOperator) o.Object {
	left := r.eval(scope, n.Left)
	if n.Operator == "??" {
		maybe := o.NewMaybe(left)
		if maybe.Ok {
			return maybe.Value
		}
		return r.eval(scope, n.Right)
	}

	if isRaise(left) {
		return left
	}

	if n.Operator == "and" {
		if !left.AsBool() {
			return o.False
		}
		return r.eval(scope, n.Right)
	}

	if n.Operator == "or" {
		if left.AsBool() {
			return o.True
		}
		return r.eval(scope, n.Right)
	}

	right := r.eval(scope, n.Right)
	if isRaise(right) {
		return right
	}

	return r.evalOperator(scope, n.Operator, left, right)
}

func (r *Evaluator) evalOperator(scope *o.Scope, op string, left, right o.Object) o.Object {
	leftTypeId := left.TypeId()
	rightTypeId := right.TypeId()

	switch {
	case op == "??":
		maybe := o.NewMaybe(left)
		if maybe.Ok {
			return maybe.Value
		}
		return right

	case op == "and":
		return o.NewBoolean(left.AsBool() && right.AsBool())

	case op == "or":
		return o.NewBoolean(left.AsBool() || right.AsBool())

	case op == "xor":
		return o.NewBoolean(left.AsBool() != right.AsBool())

	case op == "..":
		return o.NewString(left.AsString() + right.AsString())

	case leftTypeId == o.NumberId && rightTypeId == o.NumberId:
		return left.(*o.Number).OnOperator(scope, op, right)

	case leftTypeId == o.BooleanId && rightTypeId == o.BooleanId:
		return left.(*o.Boolean).OnOperator(scope, op, right)

	case leftTypeId == o.StringId && rightTypeId == o.StringId:
		return left.(*o.String).OnOperator(scope, op, right)

	case op == "==":
		return o.NewBoolean(left.Id() == right.Id())

	case op == "!=":
		return o.NewBoolean(left.Id() != right.Id())

	case leftTypeId != rightTypeId:
		return scope.Interrupt(o.Raise("types incompatible for operation '%s' '%s' '%s'", leftTypeId, op, rightTypeId))

	default:
		return scope.Interrupt(o.Raise("unknown binary operator '%s'", op))
	}
}

func (r *Evaluator) evalCall(scope *o.Scope, n *ast.Call) o.Object {
	obj := r.eval(scope, n.Target)
	if isRaise(obj) {
		return obj
	}

	// Inject `this` if it is a method call
	args := []o.Object{}
	if p := obj.Parent(); p != nil {
		args = append(args, p)
	}
	for _, arg := range n.Arguments {
		item := r.eval(scope, arg)
		if isRaise(item) {
			return item
		}
		args = append(args, item)
	}

	return r.call(scope, obj, args)
}

func (r *Evaluator) evalInstantiate(scope *o.Scope, n *ast.Instantiate) o.Object {
	target := r.eval(scope, n.Target)
	if isRaise(target) {
		return target
	}

	obj := r.instantiate(scope, target)
	if isRaise(obj) {
		return obj
	}

	for i := 0; i < len(n.Elements); i += 2 {
		key := n.Elements[i].(*ast.String)
		rawValue := n.Elements[i+1]

		value := r.eval(scope, rawValue)
		if isRaise(value) {
			return value
		}

		left := obj.GetProperty(key.Value)
		if left == nil {
			return scope.Interrupt(o.Raise("property '%s' not found in type '%s'", key.Value, target.TypeId()))
		}
		left.SetParent(obj)
		r.assign(scope, "=", key.Value, left, value)
	}

	return obj
}

func (r *Evaluator) instantiate(scope *o.Scope, target o.Object) o.Object {
	switch target := target.(type) {
	case *o.DataType:
		obj := target.Instantiate(scope)
		if isRaise(obj) {
			return obj
		}

		for name, node := range target.Attributes {
			value := r.eval(scope, node)
			if isRaise(value) {
				return value
			}

			obj.SetProperty(name, value)
		}
		return obj

	default:
		switch {
		case target.Type().TypeId() == o.TypeId:
			t := target.(o.ObjectType)
			return t.Instantiate(scope)
		default:
			return scope.Interrupt(o.Raise("type '%s' is not instantiable", target.TypeId()))
		}
	}
}

func (r *Evaluator) call(scope *o.Scope, target o.Object, args []o.Object) o.Object {
	switch target := target.(type) {
	case *o.BuiltinFunction:
		return r.callBuiltinFunction(scope, target, args)

	case *o.Function:
		return r.callFunction(scope, target, args)

	case *o.DataType:
		return r.callDataType(scope, target, args)

	default:
		if target.Type().TypeId() == o.TypeId {
			t := target.(o.ObjectType)
			return r.callType(scope, t, args)
		}
	}

	return scope.Interrupt(o.Raise("type '%s' is not callable", target.TypeId()))
}

func (r *Evaluator) callDataType(scope *o.Scope, dt *o.DataType, args []o.Object) o.Object {
	if len(args) > 0 {
		return dt.Convert(scope, args[0])
	}
	return r.instantiate(scope, dt)
}

func (r *Evaluator) callFunction(scope *o.Scope, fn *o.Function, args []o.Object) o.Object {
	fnScope := fn.Scope.New()
	params := &ast.Tuple{Elements: fn.Parameters}
	ret := r.resolveAssignment(scope, ":=", params, o.NewTuple(args...))
	if isRaise(ret) {
		return ret
	}

	ret = r.eval(fnScope, fn.Body)
	if t := asReturn(ret); t != nil {
		return t.Value
	}

	return ret
}

func (r *Evaluator) callType(scope *o.Scope, ot o.ObjectType, args []o.Object) o.Object {
	if len(args) > 0 {
		return ot.Convert(scope, args[0])
	}
	return ot.Instantiate(scope)
}

func (r *Evaluator) callBuiltinFunction(scope *o.Scope, fn *o.BuiltinFunction, args []o.Object) o.Object {
	ret := fn.Call(scope, args...)

	// Treat generator functions
	if isIteration(ret) {
		stream := ret.(*o.StreamIteration).Stream

		if stream.Finished {
			return streamFinishedError()
		}

		var ret o.Object
		if stream.Fn != nil {
			ret = r.eval(stream.Scope, stream.Fn.Body)
		} else {
			ret = stream.InternalFn(stream.Scope)
		}

		if t := asYield(ret); t != nil {
			return o.NewMaybe(t.Value)
		}

		stream.Finished = true
		return streamFinishedError()
	}

	return ret
}

func (r *Evaluator) evalAssignment(scope *o.Scope, n *ast.Assignment) o.Object {
	right := r.eval(scope, n.Right)
	if isRaise(right) {
		return right
	}

	isTuple := right.TypeId() == o.TupleId
	if !isTuple {
		right = o.NewTuple(right)
	}

	ret := r.resolveAssignment(scope, n.Operator, n.Left, right)
	if isRaise(ret) {
		return ret
	}

	if r := ret.(*o.Tuple); len(r.Elements) == 1 {
		return r.Elements[0]
	}
	return ret
}

func (r *Evaluator) resolveAssignment(scope *o.Scope, op string, left ast.Node, right o.Object) o.Object {
	switch left := left.(type) {
	case *ast.Identifier:
		return r.assign(scope, op, left.Value, scope.GetGlobal(left.Value), right)

	case *ast.Access:
		target := r.eval(scope, left)
		if isRaise(target) {
			return target
		}

		id := left.Right.(*ast.Identifier).Value
		return r.assign(scope, op, id, target, right)

	case *ast.Index:
		target := r.eval(scope, left.Target)
		if isRaise(target) {
			return target
		}

		index := r.eval(scope, left.Index)
		if isRaise(index) {
			return index
		}

		if index.TypeId() != o.TupleId {
			index = o.NewTuple(index)
		}

		return target.OnIndexAssign(scope, index.(*o.Tuple), right)

	case *ast.Tuple:
		if right.TypeId() != o.TupleId {
			right = o.NewTuple(right)
		}

		if isRaise(right) {
			return right
		}

		result := []o.Object{}
		tuple := right.(*o.Tuple)
		var ret o.Object
		j := 0
		for _, l := range left.Elements {
			s, isSpread := l.(*ast.Spread)

			if isSpread {
				var spreadAmount int
				var list o.Object
				if j > len(tuple.Elements)-1 {
					spreadAmount = 0
					list = o.NewList()

				} else {
					spreadAmount = len(tuple.Elements) - (len(left.Elements) - 1)
					from := min(j, len(tuple.Elements)-1)
					to := j + max(spreadAmount, 0)
					list = o.NewList(tuple.Elements[from:to]...)
				}

				ret = r.resolveAssignment(scope, op, s.Target, list)

				j += spreadAmount

			} else {
				if j >= len(tuple.Elements) {
					return scope.Interrupt(o.Raise("trying to unpack more elements than available in tuple assignment (expected %d, got %d)", len(left.Elements), len(tuple.Elements)))
				}

				ret = r.resolveAssignment(scope, op, l, tuple.Elements[j])
				j++
			}

			if isRaise(ret) {
				return ret
			}

			result = append(result, ret)
		}

		return o.NewTuple(result...)

	default:
		return scope.Interrupt(o.Raise("unknown node type '%v' in assignment", left))
	}
}

func (r *Evaluator) assign(scope *o.Scope, op, identifier string, left o.Object, right o.Object) o.Object {
	if op == "=" {
		switch {
		case left == nil:
			return scope.Interrupt(o.Raise("identifier '%s' is undefined. Assign a new variable using ':=' operator", identifier))

		case left.TypeId() == o.MaybeId:
			maybe := left.(*o.Maybe)
			if right.TypeId() == o.MaybeId {
				maybe.Set(right.(*o.Maybe).Result())
				return maybe
			}

			if !maybe.IsType(right) {
				return scope.Interrupt(o.Raise("cannot assign value of type '%s' to variable of type 'Maybe(%s)'", right.TypeId(), maybe.ValueType))
			}
			maybe.Set(right)
			return maybe

		case left.TypeId() != right.TypeId():
			return scope.Interrupt(o.Raise("cannot assign value of type '%s' to variable of type '%s'. If you want to change types, reassign the variable with ':=' operator", right.TypeId(), left.TypeId()))
		}

	} else if op == ":=" {
		if left != nil && left.Parent() != nil {
			return scope.Interrupt(o.Raise("cannot reassign property '%s' of a type '%s'", identifier, left.TypeId()))
		}
	}

	if identifier == "_" {
		return right
	}

	if op == "=" {
		if p := left.Parent(); p != nil {
			p.SetProperty(identifier, right)
		} else {
			scope.SetGlobal(identifier, right)
		}
	} else {
		scope.SetLocal(identifier, right)
	}

	return right
}

func (r *Evaluator) evalSpread(scope *o.Scope, n *ast.Spread) o.Object {
	if n.In {
		return scope.Interrupt(o.Raise("spread in operator '...' is not supported in this context"))
	}

	target := r.eval(scope, n.Target)
	if isRaise(target) {
		return target
	}

	tuple := o.TupleTypeObj.Convert(scope, target)

	return tuple
}

func (r *Evaluator) evalAccess(scope *o.Scope, n *ast.Access) o.Object {
	left := r.eval(scope, n.Left)
	if isRaise(left) {
		return left
	}

	switch right := n.Right.(type) {
	case *ast.Identifier:
		p := left.GetProperty(right.Value)
		if p == nil {
			return scope.Interrupt(o.Raise("property '%s' not found in type '%s'", right.Value, left.TypeId()))
		}

		if left.Type() != o.TypeTypeObj {
			data, ok := left.Type().(*o.DataType)
			if ok && data.Attributes[right.Value] != nil {
				return p
			}
			p.SetParent(left)
		}
		return p

	default:
		return scope.Interrupt(o.Raise("unknown node type '%v' in access", n.Right))
	}
}

func (r *Evaluator) evalIndex(scope *o.Scope, n *ast.Index) o.Object {
	target := r.eval(scope, n.Target)
	if isRaise(target) {
		return target
	}

	index := r.eval(scope, n.Index)
	if isRaise(index) {
		return index
	}

	if index.TypeId() != o.TupleId {
		index = o.NewTuple(index)
	}

	return target.OnIndex(scope, index.(*o.Tuple))
}

func (r *Evaluator) evalWrap(scope *o.Scope, n *ast.Wrap) o.Object {
	target := r.eval(scope, n.Target)

	if t := asRaise(target); t != nil {
		return o.NewMaybe(t.Value)
	}

	return o.NewMaybe(target)
}

func (r *Evaluator) evalUnwrap(scope *o.Scope, n *ast.Unwrap) o.Object {
	target := r.eval(scope, n.Target)
	if isRaise(target) {
		return target
	}

	maybe := o.NewMaybe(target)
	return o.NewTuple(maybe.Error, maybe.Value)
}

// ----------------------------------------------------------------------------
// Control Flow Evaluation
// ----------------------------------------------------------------------------
func (r *Evaluator) evalBlock(scope *o.Scope, n *ast.Block) o.Object {
	const BlockReturnKey = "$block-return"

	// TODO Save last statements value in the scope because it can be the last result
	var blockScope *o.Scope
	var curStatement = 0

	ar := scope.ActiveRecord()
	if ar != nil {
		state := ar.(*BlockRecord)
		blockScope = state.Scope
		curStatement = state.Statement
	} else {
		blockScope = scope.New()
		blockScope.SetLocal(BlockReturnKey, o.False)
	}
	scope.SetActiveRecord(nil)

	for i := curStatement; i < len(n.Expressions); i++ {
		statement := n.Expressions[i]
		result := r.eval(blockScope, statement)

		if t := asYield(result); t != nil {
			idx := i
			if t.TriggeredScope == blockScope {
				idx = i + 1
			}

			scope.SetActiveRecord(&BlockRecord{
				Scope:     blockScope,
				Statement: idx,
			})
		}

		if t := asInterruption(result); t != nil {
			if t.Category != o.BreakId && t.Category != o.ContinueId {
				blockScope.SetLocal(BlockReturnKey, o.False)
			}

			return result
		}

		blockScope.SetLocal(BlockReturnKey, result)
	}

	return blockScope.GetLocal(BlockReturnKey)
}

func (r *Evaluator) evalReturn(scope *o.Scope, n *ast.Return) o.Object {
	right := r.eval(scope, n.Expression)
	if isRaise(right) {
		return right
	}

	return scope.Interrupt(o.ReturnWith(right))
}

func (r *Evaluator) evalRaise(scope *o.Scope, n *ast.Raise) o.Object {
	right := r.eval(scope, n.Expression)
	if isRaise(right) {
		return right
	}

	return scope.Interrupt(o.RaiseWith(right))
}

func (r *Evaluator) evalYield(scope *o.Scope, n *ast.Yield) o.Object {
	if n.Break {
		return scope.Interrupt(o.ReturnWith(o.False))
	}

	right := r.eval(scope, n.Expression)
	if isRaise(right) {
		return right
	}

	return scope.Interrupt(o.YieldWith(right))
}

func (r *Evaluator) evalBreak(scope *o.Scope, _ *ast.Break) o.Object {
	return scope.Interrupt(&o.Interruption{
		Category: o.BreakId,
	})
}

func (r *Evaluator) evalContinue(scope *o.Scope, _ *ast.Continue) o.Object {
	return scope.Interrupt(&o.Interruption{
		Category: o.ContinueId,
	})
}

func (r *Evaluator) evalIf(scope *o.Scope, n *ast.If) o.Object {
	var ifScope *o.Scope
	var condition *bool
	ar := scope.ActiveRecord()
	if ar != nil {
		state := ar.(*IfRecord)
		condition = &state.Condition
		ifScope = state.Scope
	} else {
		ifScope = scope.New()
	}
	scope.SetActiveRecord(nil)

	if condition == nil {
		var res o.Object = o.False
		for _, condition := range n.Conditions {
			res = r.eval(ifScope, condition)
			if isRaise(res) {
				return res
			}
		}

		b := res.AsBool() // use last one
		condition = &b
	}

	var ret o.Object
	if *condition {
		ret = r.eval(ifScope, n.TrueExpression)
	} else {
		ret = r.eval(ifScope, n.FalseExpression)
	}

	if t := asYield(ret); t != nil {
		scope.SetActiveRecord(&IfRecord{
			Scope:     ifScope,
			Condition: *condition,
		})
	}

	return ret
}

func (r *Evaluator) evalFor(scope *o.Scope, n *ast.For) o.Object {
	var forScope *o.Scope
	conditionSolved := false
	ar := scope.ActiveRecord()
	if ar != nil {
		state := ar.(*ForRecord)
		forScope = state.Scope
		conditionSolved = true
	} else {
		conditionSolved = false
		forScope = scope.New()
		forScope.SetLocal(ForReturnKey, o.False)
	}
	scope.SetActiveRecord(nil)

	for {
		if !conditionSolved && len(n.Conditions) > 0 {
			ret, shouldReturn := r.checkForCondition(forScope, n)
			if shouldReturn {
				return ret
			}
		}

		res := r.eval(forScope, n.Expression)
		if ret := asYield(res); ret != nil {
			scope.SetActiveRecord(&ForRecord{
				Scope: forScope,
			})
		}

		conditionSolved = false
		if t := asInterruption(res); t != nil {
			if t.Category == o.BreakId {
				break
			}
			if t.Category == o.ContinueId {
				continue
			}
			forScope.SetLocal(ForReturnKey, t.Value)
			return t
		}

		forScope.SetLocal(ForReturnKey, res)
	}

	return forScope.GetLocal(ForReturnKey)
}

func (r *Evaluator) checkForCondition(scope *o.Scope, n *ast.For) (res o.Object, shouldReturn bool) {
	// Run first conditions
	for _, condition := range n.Conditions[:len(n.Conditions)-1] {
		res = r.eval(scope, condition)
		if isRaise(res) {
			return res, true
		}
	}

	// Perform the check for the last condition
	last := n.Conditions[len(n.Conditions)-1]

	// ... no 'in' expression, just evaluate
	if n.InExpression == nil {
		res = r.eval(scope, last)
		if isRaise(res) {
			return res, true
		}

		if !res.AsBool() {
			return scope.GetLocal(ForReturnKey), true
		}
		return res, false
	}

	// ... 'in' expression, evaluate the expression and check if it is finished
	// ... ... first time, initialize the iterator
	stream := scope.GetLocal(ForInKey)
	if stream == nil {
		res = r.eval(scope, n.InExpression)
		if isRaise(res) {
			return res, true
		}

		stream = o.StreamTypeObj.Convert(scope, res)
		if isRaise(stream) {
			return stream, true
		}

		scope.SetLocal(ForInKey, stream)
	}

	// ... ... resolve the iterator
	iter := r.callBuiltinFunction(scope, o.Stream_Next, []o.Object{stream})
	if isRaise(iter) {
		return iter, true
	}

	// assign maybe into variables
	maybe := iter.(*o.Maybe)
	if !maybe.Ok {
		return scope.GetLocal(ForReturnKey), true
	}

	if _, ok := last.(*ast.Tuple); !ok {
		last = &ast.Tuple{Token: last.GetToken(), Elements: []ast.Node{last}}
	}

	ret := r.resolveAssignment(scope, ":=", last, maybe.Value)
	if isRaise(ret) {
		return ret, true
	}
	return ret, false
}

func (r *Evaluator) evalWith(scope *o.Scope, n *ast.With) o.Object {
	var withScope *o.Scope
	ar := scope.ActiveRecord()
	if ar != nil {
		state := ar.(*WithRecord)
		withScope = state.Scope
	} else {
		withScope = scope.New()
	}
	scope.SetActiveRecord(nil)

	if ar == nil {
		var res o.Object = o.False
		res = r.eval(withScope, n.Condition)
		if isRaise(res) {
			return res
		}
	}

	ret := r.eval(withScope, n.Expression)
	if t := asYield(ret); t != nil {
		scope.SetActiveRecord(&WithRecord{
			Scope: withScope,
		})
	}

	return ret
}

func (r *Evaluator) evalMatch(scope *o.Scope, n *ast.Match) o.Object {
	var matchScope *o.Scope
	var caseScope *o.Scope
	var caseIdx int = -1
	ar := scope.ActiveRecord()
	if ar != nil {
		state := ar.(*MatchRecord)
		caseIdx = state.Case
		matchScope = state.Scope
		caseScope = state.CaseScope
	} else {
		matchScope = scope.New()
	}
	scope.SetActiveRecord(nil)

	if caseIdx == -1 {
		expression := r.eval(matchScope, n.Expression)
		if isRaise(expression) {
			return expression
		}

		for i := 0; i < len(n.Cases); i += 2 {
			condition := n.Cases[i]

			caseScope = matchScope.New()
			if r.match(matchScope, condition, expression) {
				caseIdx = i
				break
			}
		}
	}

	if caseIdx == -1 {
		return o.False
	}

	ret := r.eval(caseScope, n.Cases[caseIdx+1])
	if t := asYield(ret); t != nil {
		scope.SetActiveRecord(&MatchRecord{
			Scope: matchScope,
			Case:  caseIdx,
		})
	}

	return ret
}

func (r *Evaluator) match(scope *o.Scope, a ast.Node, b o.Object) bool {
	switch a := a.(type) {
	case *ast.Tuple:
		if b.TypeId() != o.TupleId {
			return false
		}

		tuple := b.(*o.Tuple)
		if len(a.Elements) != len(tuple.Elements) {
			return false
		}

		for i, e := range a.Elements {
			if !r.match(scope, e, tuple.Elements[i]) {
				return false
			}
		}

		return true

	case *ast.Assignment:
		r.resolveAssignment(scope, ":=", a.Left, b)
		return true

	default:
		if a.GetToken().IsLiteral("_") {
			return true
		}

		left := r.eval(scope, a)
		if isRaise(left) {
			return false
		}

		return r.evalOperator(scope, "==", left, b).AsBool()
	}
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func asInterruption(obj o.Object) *o.Interruption {
	intr, ok := obj.(*o.Interruption)
	if ok {
		return intr
	}
	return nil
}

func isRaise(obj o.Object) bool {
	intr, ok := obj.(*o.Interruption)
	return ok && intr.Category == o.RaiseId
}

func asRaise(obj o.Object) *o.Interruption {
	intr, ok := obj.(*o.Interruption)
	if ok && intr.Category == o.RaiseId {
		return intr
	}
	return nil
}

func asReturn(obj o.Object) *o.Interruption {
	intr, ok := obj.(*o.Interruption)
	if ok && intr.Category == o.ReturnId {
		return intr
	}
	return nil
}

func asYield(obj o.Object) *o.Interruption {
	intr, ok := obj.(*o.Interruption)
	if ok && intr.Category == o.YieldId {
		return intr
	}
	return nil
}

func isIteration(obj o.Object) bool {
	_, ok := obj.(*o.StreamIteration)
	return ok
}

func streamFinishedError() o.Object {
	return o.NewMaybe(o.NewErrorFromString("Stream finished"))
}
