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
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

		// expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpr(node.Operator, right)

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

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, s := range statements {
		result = Eval(s)
	}

	return result
}

func boolToBoolean(b bool) *object.Boolean {
	if b {
		return ConstTrue
	}
	return ConstFalse
}

func evalPrefixExpr(op token.Type, right object.Object) object.Object {
	switch op {

	}
	return nil
}
