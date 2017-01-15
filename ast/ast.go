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

// IfExpression ::= 'if' expression ('{' | '->') BlockStatement '}'?
type IfExpression struct {
	Token token.Token // token.If
	Cond  Expression
	Do    *BlockStatement
	Else  *BlockStatement
}

// FunctionLiteral ::= '|' (Identifier | (Identifier ',')?)* ('{' | '->') BlockStatement '}'?
type FunctionLiteral struct {
	Token  token.Token // The first '|' bar token
	Params []*Identifier
	Body   *BlockStatement
}

// BlockStatement ::= Statement*
type BlockStatement struct {
	Token      token.Token // { or ->
	Statements []Statement
}

// CallExpression ::= Identifier ('|' (Expression | (Expression ',')?)* '|' | '!')
type CallExpression struct {
	Token token.Token // the | or ! token
	Func  Expression  // either an Identifier or function literal
	Args  []Expression
}

// ExpressionStatement ::= (IntegerLiteral | FloatLiteral | StringLiteral | Operator | Identifier) Expression?
type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

// PrefixExpression ::= Operator Expression
type PrefixExpression struct {
	Token    token.Token // prefix token ! & -
	Operator token.Type
	Right    Expression
}

// InfixExpression ::= Expression Operator Expression
type InfixExpression struct {
	Token    token.Token // The operator
	Left     Expression
	Operator token.Type
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

// BooleanLiteral ::= True | False
type BooleanLiteral struct {
	// token.True | token.False
	Token token.Token
	Value bool
}
