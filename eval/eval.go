package eval

import (
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

// Eval evaluates the program node and returns an object as a result
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// statements
	case *ast.Program:
		return evalProgram(node)
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Eval(node.Value)}
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

		// expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpr(node.Operator, left, right)
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

		if result.Type() == object.ReturnType {
			return result.(*object.ReturnValue).Value
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, s := range block.Statements {
		result = Eval(s)

		if result != nil && result.Type() == object.ReturnType {
			return result
		}
	}

	return result
}

func evalIfExpr(node *ast.IfExpression) object.Object {
	cond := Eval(node.Cond)

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

func evalPrefixExpr(op token.Type, right object.Object) object.Object {
	switch op {
	case token.Bang:
		return evalBangOperatorExpr(right)
	case token.Minus:
		return evalMinusPrefixOperatorExpr(right)
	}
	return nil
}

func evalInfixExpr(op token.Type, left object.Object, right object.Object) object.Object {
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
	if op == token.Equal {
		return boolToBoolean(left == right)
	} else if op == token.NotEqual {
		return boolToBoolean(left != right)
	}

	// otherwise 2 objects that don't match
	return nil
}

func evalIntegerInfixExpr(op token.Type, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
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
		return ConstNil
	}
}

func evalFloatInfixExpr(op token.Type, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch op {
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
		return ConstNil
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
		return ConstFalse
	}
}

func evalMinusPrefixOperatorExpr(right object.Object) object.Object {
	switch right.Type() {
	case object.IntType:
		v := right.(*object.Integer).Value
		return &object.Integer{Value: -v}
	case object.FloatType:
		v := right.(*object.Float).Value
		return &object.Float{Value: -v}
	default:
		return ConstNil
	}
}
