// ast/ast.go
// contains structs for ast
// self documents ebnf or lang
// interface implementations are split into ast/ast_interface.go
// to increase readability or tree structure

package ast

import "jacob/black/token"

// Program ::= Statement*
type Program struct {
	Statements []Statement
}

// LetStatement ::= 'let' Identifier '=' Expression
type LetStatement struct {
	Token token.Token // token.Let
	Name  *Identifier
	Value Expression
}

// ReturnStatement ::= 'ret' expression
type ReturnStatement struct {
	Token token.Token // token.Return
	Value Expression
}

// ExpressionStatement ::= (IntegerLiteral | FloatLiteral | StringLiteral | Operator | Identifier) Expression?
type ExpressionStatement struct {
	// first token of the expression
	Token      token.Token
	Expression Expression
}

// PrefixExpression ::= Operator Expression
type PrefixExpression struct {
	// prefix token ! & -
	Token    token.Token
	Operator string
	Right    Expression
}

// InfixExpression ::= Expression Operator Expression
type InfixExpression struct {
	// The operator
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

// Identifier ::= name
type Identifier struct {
	Token token.Token // token.Identifier
	Value string
}

// IntegerLiteral ::= int of 64 bits
type IntegerLiteral struct {
	Token token.Token // token.Int
	Value int64
}

// FloatLiteral ::= float of 64 bits
type FloatLiteral struct {
	Token token.Token // token.Float
	Value float64
}
