package eval

import (
	"fmt"
	"jacob/dusk/ast"
	"jacob/dusk/object"
	"jacob/dusk/token"
	"math"
)

var (
	// ConstTrue is the only and only true value
	ConstTrue = &object.Boolean{Value: true}
	// ConstFalse is the only and only false value
	ConstFalse = &object.Boolean{Value: false}
	// ConstNil is the only and only nil
	ConstNil = &object.Nil{}
)

func newError(pos token.Position, format string, v ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, v...), Pos: pos}
}

func isError(o object.Object) bool {
	if o != nil {
		return o.Type() == object.ErrorType
	}
	return false
}

func bottomEnv(id *ast.AccessIdentifier, env *object.Environment) (*object.Environment, string, object.Object) {
	currentEnv := env

	for i, v := range id.Values {
		if i < len(id.Values)-1 {
			if val, ok := currentEnv.Get(v); ok {
				if val.Type() != object.FunctionType {
					return nil, v, newError(id.Token.Pos, "cannot use '.' on type '%s'. Must be function", val.Type())
				}

				val := val.(*object.Function)
				currentEnv = val.Env
			} else {
				return nil, v, newError(id.Token.Pos, "identifier not found in context: %s", v)
			}
		}
	}

	return currentEnv, id.Values[len(id.Values)-1], nil
}

// Eval evaluates the program node and returns an object as a result
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

		// expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Token, right)
	case *ast.InfixExpression:
		if node.Operator == token.Assign {
			return evalAssign(node, env)
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpr(node.Token, left, right)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpr(node.Token, left, index)
	case *ast.IfExpression:
		return evalIfExpr(node, env)
	case *ast.WhileExpression:
		return evalWhileExpr(node, env)
	case *ast.CallExpression:
		function := Eval(node.Func, env)
		if isError(function) {
			return function
		}

		args, err := evalExpressions(node.Args, env)
		if err != nil {
			return err
		}

		return doFunction(node.Token, function, args)

		// literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.BooleanLiteral:
		return boolToBoolean(node.Value)
	case *ast.NilLiteral:
		return ConstNil
	case *ast.FunctionLiteral:
		return &object.Function{Params: node.Params, Body: node.Body, Env: env}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elems, err := evalExpressions(node.Elements, env)
		if len(elems) == 1 && err != nil {
			return err
		}
		return &object.Array{Elements: elems}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.AccessIdentifier:
		return evalAccessIdentifier(node, env)
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, s := range program.Statements {
		result = Eval(s, env)

		if result != nil {
			// pass up the return type to the top level
			switch result.Type() {
			case object.ReturnType:
				// unwrap the return value
				return result.(*object.ReturnValue).Value
			case object.ErrorType:
				return result.(*object.Error)
			}
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, s := range block.Statements {
		result = Eval(s, env)

		if result != nil {
			if result.Type() == object.ReturnType || result.Type() == object.ErrorType {
				return result
			}
		}
	}

	return result
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) ([]object.Object, object.Object) {
	var evaluated []object.Object

	for _, e := range expressions {
		evaled := Eval(e, env)
		if isError(evaled) {
			return []object.Object{}, evaled
		}
		evaluated = append(evaluated, evaled)
	}

	return evaluated, nil
}

func evalIfExpr(node *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(node.Cond, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(node.Do, env)
	} else if node.Else != nil {
		return Eval(node.Else, env)
	}

	return ConstNil
}

func evalWhileExpr(node *ast.WhileExpression, env *object.Environment) object.Object {
	//var lastEval object.Object = ConstNil

	for {
		cond := Eval(node.Cond, env)
		if isError(cond) {
			return cond
		}

		// while should return nothing
		if !isTruthy(cond) {
			return nil
		}

		Eval(node.Do, env)

		if node.Then != nil {
			then := Eval(node.Then, env)
			if isError(then) {
				return then
			}
		}
	}
}

// isTruthy - everything is true execpt for false and nil
func isTruthy(o object.Object) bool {
	switch o {
	case ConstFalse, ConstNil:
		return false
	default:
		// special case: 0 or 0.0 is not truthy
		switch o.Type() {
		case object.IntType:
			if o.(*object.Integer).Value == 0 {
				return false
			}
		case object.FloatType:
			if o.(*object.Float).Value == 0.0 {
				return false
			}
		}

		return true
	}
}

func boolToBoolean(b bool) *object.Boolean {
	if b {
		return ConstTrue
	}
	return ConstFalse
}

func evalPrefixExpr(op token.Token, right object.Object) object.Object {
	switch op.Type {
	case token.Bang:
		return evalBangOperatorExpr(right)
	case token.Minus:
		return evalMinusPrefixOperatorExpr(op.Pos, right)
	default:
		return newError(op.Pos, "unknown operator '%s' for type '%s'", op, right.Type())
	}
}

func evalInfixExpr(op token.Token, left, right object.Object) object.Object {

	// catch early type errors
	if !left.CanApply(op.Type, right.Type()) && (op.Type != token.Equal && op.Type != token.NotEqual) {
		return newError(op.Pos, "cannot apply operator '%s' for type '%s' and '%s'", op, left.Type(), right.Type())
	}

	// test and convert int to float if needed
	if left.Type() == object.IntType && right.Type() == object.IntType {
		// both ints. easy
		return evalIntegerInfixExpr(op, left, right)
	} else if left.Type() == object.FloatType && right.Type() == object.FloatType {
		// both floats. easy
		return evalFloatInfixExpr(op, left, right)
	}

	// one of them must be a int and the other a float
	if left.Type() == object.IntType || right.Type() == object.IntType {
		if left.Type() == object.FloatType {
			// left is float, right is int.
			// premote right to float
			val := right.(*object.Integer).Value
			right = &object.Float{Value: float64(val)}
			return evalFloatInfixExpr(op, left, right)

		} else if right.Type() == object.FloatType {
			// right is float, left is int
			// premote left to float
			val := left.(*object.Integer).Value
			left = &object.Float{Value: float64(val)}

			return evalFloatInfixExpr(op, left, right)
		}
	}

	// two strings
	if left.Type() == object.StringType && right.Type() == object.StringType {
		return evalStringInfixExpr(op, left, right)
	}

	// two arrays
	if left.Type() == object.ArrayType && right.Type() == object.ArrayType {
		return evalArrayInfixExpr(op, left, right)
	}

	// compare actual runtime object
	if op.Type == token.Equal {
		return boolToBoolean(left == right)
	} else if op.Type == token.NotEqual {
		return boolToBoolean(left != right)
	}

	// otherwise 2 objects that don't match
	return newError(op.Pos, "unknown operator '%s' for type '%s' and '%s'", op, left.Type(), right.Type())
}

func evalIndexExpr(op token.Token, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayType && index.Type() == object.IntType:
		return evalArrayIndexExpr(op.Pos, left, index)
	case left.Type() == object.StringType && index.Type() == object.IntType:
		return evalStringIndexExpr(op.Pos, left, index)
	default:
		return newError(op.Pos, "index operator not supported on type '%s'", left.Type())
	}
}

func evalArrayIndexExpr(pos token.Position, arr, index object.Object) object.Object {
	array := arr.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(array.Elements) - 1)

	if i < 0 {
		i = max + i + 1
	}

	if i < 0 || i > max {
		return newError(pos, "index '%d' out of bounds of array. Max '%d'", i, max)
	}

	return array.Elements[i]
}

func evalStringIndexExpr(pos token.Position, str, index object.Object) object.Object {
	s := str.(*object.String)
	i := index.(*object.Integer).Value
	max := int64(len(s.Value) - 1)

	if i < 0 {
		i = max + i + 1
	}

	if i < 0 || i > max {
		return newError(pos, "index '%d' out of bounds of string. Max '%d'", i, max)
	}

	return &object.String{Value: string(s.Value[i])}
}

func evalStringInfixExpr(op token.Token, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch op.Type {
	case token.Plus:
		return &object.String{Value: leftVal + rightVal}
	case token.Equal:
		return boolToBoolean(leftVal == rightVal)
	case token.NotEqual:
		return boolToBoolean(leftVal == rightVal)
	default:
		return newError(op.Pos, "unknown operator '%s' for type '%s' and '%s'", op.Type, left.Type(), right.Type())
	}
}

func evalArrayInfixExpr(op token.Token, left, right object.Object) object.Object {
	leftVals := left.(*object.Array).Elements
	rightVals := right.(*object.Array).Elements

	switch op.Type {
	case token.Plus:
		return &object.Array{Elements: append(leftVals, rightVals...)}
	case token.Equal:
		if len(leftVals) != len(rightVals) {
			return ConstFalse
		}

		for i := range leftVals {
			if evalInfixExpr(op, leftVals[i], rightVals[i]) == ConstFalse {
				return ConstFalse
			}
		}
		return ConstTrue
	case token.NotEqual:
		if len(leftVals) != len(rightVals) {
			return ConstTrue
		}

		for i := range leftVals {
			if evalInfixExpr(op, leftVals[i], rightVals[i]) == ConstFalse {
				return ConstTrue
			}
		}
		return ConstFalse
	default:
		return newError(op.Pos, "unknown operator '%s' for type '%s' and '%s'", op.Type, left.Type(), right.Type())
	}
}

func evalIntegerInfixExpr(op token.Token, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.Plus:
		return &object.Integer{Value: leftVal + rightVal}
	case token.Minus:
		return &object.Integer{Value: leftVal - rightVal}
	case token.Times:
		return &object.Integer{Value: leftVal * rightVal}
	case token.Divide:
		if rightVal == 0 {
			return newError(op.Pos, "cannot divide %d by 0", leftVal)
		}
		return &object.Integer{Value: leftVal / rightVal}
	case token.Exp:
		return &object.Integer{Value: int64(math.Pow(float64(leftVal), float64(rightVal)))}
	case token.Mod:
		if rightVal == 0 {
			return newError(op.Pos, "cannot modulo %d by 0", leftVal)
		}
		return &object.Integer{Value: leftVal % rightVal}
	case token.Less:
		return boolToBoolean(leftVal < rightVal)
	case token.Greater:
		return boolToBoolean(leftVal > rightVal)
	case token.Equal:
		return boolToBoolean(leftVal == rightVal)
	case token.NotEqual:
		return boolToBoolean(leftVal != rightVal)
	default:
		return newError(op.Pos, "unknown operator '%s' for type '%s' and '%s'", op.Type, left.Type(), right.Type())
	}
}

func evalFloatInfixExpr(op token.Token, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch op.Type {
	case token.Plus:
		return &object.Float{Value: leftVal + rightVal}
	case token.Minus:
		return &object.Float{Value: leftVal - rightVal}
	case token.Times:
		return &object.Float{Value: leftVal * rightVal}
	case token.Divide:
		if rightVal == 0 {
			return newError(op.Pos, "cannot divide %d by 0", leftVal)
		}
		return &object.Float{Value: leftVal / rightVal}
	case token.Exp:
		return &object.Float{Value: math.Pow(leftVal, rightVal)}
	case token.Mod:
		if rightVal == 0 {
			return newError(op.Pos, "cannot modulo %d by 0", leftVal)
		}
		return &object.Float{Value: math.Mod(leftVal, rightVal)}
	case token.Less:
		return boolToBoolean(leftVal < rightVal)
	case token.Greater:
		return boolToBoolean(leftVal > rightVal)
	case token.Equal:
		return boolToBoolean(leftVal == rightVal)
	case token.NotEqual:
		return boolToBoolean(leftVal != rightVal)
	default:
		return newError(op.Pos, "unknown operator '%s' for type '%s' and '%s'", op.Type, left.Type(), right.Type())
	}
}

func evalBangOperatorExpr(right object.Object) object.Object {
	switch right {
	case ConstTrue:
		return ConstFalse
	case ConstFalse:
		return ConstTrue
	case ConstNil:
		return ConstTrue
	default:
		return boolToBoolean(!isTruthy(right))
	}
}

func evalMinusPrefixOperatorExpr(pos token.Position, right object.Object) object.Object {
	switch right.Type() {
	case object.IntType:
		v := right.(*object.Integer).Value
		return &object.Integer{Value: -v}
	case object.FloatType:
		v := right.(*object.Float).Value
		return &object.Float{Value: -v}
	default:
		return newError(pos, "unknown operator '-' for type '%s'", right.Type())
	}
}

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(id.Value); ok {
		return val
	}

	if builtin, ok := builtins[id.Value]; ok {
		return builtin
	}

	return newError(id.Token.Pos, "identifier not found: %s", id.Value)
}

func evalAccessIdentifier(id *ast.AccessIdentifier, env *object.Environment) object.Object {

	bottom, last, err := bottomEnv(id, env)
	if err != nil {
		return err
	}

	if val, ok := bottom.Get(last); ok {
		return val
	}

	if builtin, ok := builtins[last]; ok {
		return builtin
	}

	return newError(id.Token.Pos, "identifier not found in context: %s", last)
}

// special case = assign operator
func evalAssign(node *ast.InfixExpression, env *object.Environment) object.Object {

	var bottom *object.Environment
	var id string

	switch l := node.Left.(type) {
	case *ast.Identifier:
		id = l.Value
		bottom = env
	case *ast.AccessIdentifier:
		b, last, err := bottomEnv(l, env)
		if err != nil {
			return err
		}

		id = last
		bottom = b
	case *ast.IndexExpression:
		array := Eval(l.Left, env)
		if isError(array) {
			return array
		}
		if array.Type() != object.ArrayType {
			return newError(l.Token.Pos, "index operator assign not supported on type '%s'", array.Type())
		}
		a := array.(*object.Array)

		index := Eval(l.Index, env)
		if isError(index) {
			return index
		}
		if index.Type() != object.IntType {
			return newError(l.Token.Pos, "index must be type 'int'. got type '%s'", index.Type())
		}
		i := index.(*object.Integer).Value

		max := int64(len(a.Elements) - 1)
		if i < 0 {
			i = max + i + 1
		}

		if i < 0 || i > max {
			return newError(l.Token.Pos, "index '%d' out of bounds of array. Max '%d'", i, max)
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		if array == right {
			return newError(l.Token.Pos, "cannot assign index of array to self")
		}

		array.(*object.Array).Elements[index.(*object.Integer).Value] = right
		return right

	default:
		return newError(node.Token.Pos, "cannot bind a literal to a value")
	}

	// check if exists already
	if val, ok := bottom.Get(id); ok {
		if isError(val) {
			return val
		}

		// eval rhs
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		// must be same type
		if val.Type() == right.Type() || (val.Type() == object.NilType || right.Type() == object.NilType) {
			v, ok := bottom.Assign(id, right)
			if ok {
				return v
			}
			return newError(node.Token.Pos, "cannot assign value to variable '%s' that does not exist", id)
		}

		return newError(node.Token.Pos, "cannot assign variable '%s' of type '%s' to value '%s' of type '%s'", id, val.Type(), right, right.Type())
	}
	return newError(node.Token.Pos, "cannot assign value to variable '%s' that does not exist", id)
}

func doFunction(t token.Token, f object.Object, args []object.Object) object.Object {

	switch function := f.(type) {
	case *object.Function:
		if len(function.Params) != len(args) {
			return newError(t.Pos, "invalid number of arguments for function. Expected %d got %d", len(function.Params), len(args))
		}

		childEnv := adoptFunctionEnv(function, args)
		evaluated := Eval(function.Body, childEnv)

		if val, ok := evaluated.(*object.ReturnValue); ok {
			return val.Value
		}

		return evaluated
	case *object.Builtin:
		return function.Fn(args...)
	default:
		return newError(t.Pos, "type '%s' not a function", f.Type())
	}
}

func adoptFunctionEnv(f *object.Function, args []object.Object) *object.Environment {
	env := object.NewChildEnvironment(f.Env)

	for i, p := range f.Params {
		env.Set(p.Value, args[i])
	}

	return env
}
