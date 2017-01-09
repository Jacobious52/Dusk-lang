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
	// token.Let
	Token token.Token
	Name  *Identifier
	Value Expression
}

// Identifier ::= name
type Identifier struct {
	// token.Identifier
	Token token.Token
	Value string
}

// ReturnStatement ::= 'ret' expression
type ReturnStatement struct {
	// token.Return
	Token token.Token
	Value Expression
}

// ExpressionStatement ::= (Number | Operator | Identifier) Expression?
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}
