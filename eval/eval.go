package eval

import (
	"fmt"
	"jacob/black/ast"
	"jacob/black/object"
	"jacob/black/token"
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

// Eval evaluates the program node and returns an object as a result
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// statements
	case *ast.Program:
		return evalProgram(node)
	case *ast.ReturnStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

		// expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Token, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpr(node.Token, left, right)
	case *ast.IfExpression:
		return evalIfExpr(node)

		// literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.BooleanLiteral:
		return boolToBoolean(node.Value)
	}
	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, s := range program.Statements {
		result = Eval(s)

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

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, s := range block.Statements {
		result = Eval(s)

		if result != nil {
			if result.Type() == object.ReturnType || result.Type() == object.ErrorType {
				return result
			}
		}
	}

	return result
}

func evalIfExpr(node *ast.IfExpression) object.Object {
	cond := Eval(node.Cond)

	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(node.Do)
	} else if node.Else != nil {
		return Eval(node.Else)
	}

	return ConstNil
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

func evalInfixExpr(op token.Token, left object.Object, right object.Object) object.Object {

	if !left.CanApply(op.Type, right.Type()) {
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

	// compare actual runtime object
	if op.Type == token.Equal {
		return boolToBoolean(left == right)
	} else if op.Type == token.NotEqual {
		return boolToBoolean(left != right)
	}

	// otherwise 2 objects that don't match
	return newError(op.Pos, "unknown operator '%s' for type '%s' and '%s'", op, left.Type(), right.Type())
}

func evalIntegerInfixExpr(op token.Token, left object.Object, right object.Object) object.Object {
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
		return &object.Integer{Value: leftVal / rightVal}
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

func evalFloatInfixExpr(op token.Token, left object.Object, right object.Object) object.Object {
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
		return &object.Float{Value: leftVal / rightVal}
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
