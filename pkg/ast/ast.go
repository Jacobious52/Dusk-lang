// ast/ast.go
// contains structs for ast
// self documents ebnf or lang
// interface implementations are split into ast/ast_interface.go
// to increase readability or tree structure

package ast

import "jacob/dusk/pkg/token"

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

// IfExpression ::= 'if' expression ('{' | ':') BlockStatement '}'? 'else' ('{' | ':')? BlockStatement '}'?
type IfExpression struct {
	Token token.Token // token.If
	Cond  Expression
	Do    *BlockStatement
	Else  *BlockStatement
}

// WhileExpression ::= 'while' expression ('{' | ':') BlockStatement '}'? 'else' ('{' | ':')? BlockStatement '}'?
type WhileExpression struct {
	Token token.Token // token.While
	Cond  Expression
	Then  Expression
	Do    *BlockStatement
}

// FunctionLiteral ::= '|' (Identifier | (Identifier ',')?)* ('{' | ':')? BlockStatement '}'?
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

// ExpressionStatement ::= (IntegerLiteral | FloatLiteral | StringLiteral | Operator | Identifier | PrefixExpression | InfixExpression | IndexExpression) Expression?
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

// AccessIdentifier ::= name '.' name ('.' name)*
type AccessIdentifier struct {
	Token  token.Token // the root (first) token.Identifier
	Values []string
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
	Token token.Token // token.True | token.False
	Value bool
}

// NilLiteral ::= Nil
type NilLiteral struct {
	Token token.Token // token.Nil
}

// StringLiteral ::= "(a...z)"
type StringLiteral struct {
	Token token.Token // token.String
	Value string
}

// ArrayLiteral ::= '[' (Expression ',')* ']'
type ArrayLiteral struct {
	Token    token.Token // token.LBracket
	Elements []Expression
}

// IndexExpression ::= Expression '[' Expression ']'
type IndexExpression struct {
	Token token.Token // token.LBracket
	Left  Expression
	Index Expression
}
