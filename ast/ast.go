package ast

import "jacob/black/token"

// Node is the the basis element of the ast
type Node interface {
	TokenLiteral() string
}

// Statement is the basis for a statment in the ast
type Statement interface {
	Node
	statementNode()
}

// Expression is the basis for a expression in the ast
type Expression interface {
	Node
	expressionNode()
}

// Program represents the whole runable program. Is the root node
type Program struct {
	Statements []Statement
}

// TokenLiteral impl for Program
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// LetStatement is the main way of declaring ans assigning objects
type LetStatement struct {
	// token.Let
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral impl for LetStatement
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Identifier is any named token.Identifier
type Identifier struct {
	// token.Identifier
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral impl for Identifier
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
