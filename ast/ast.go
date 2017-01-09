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

// ExpressionStatement ::= (Number | Literal | Operator | Identifier) Expression?
type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

// Identifier ::= name
type Identifier struct {
	// token.Identifier
	Token token.Token
	Value string
}
