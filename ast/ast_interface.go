package ast

import "bytes"

// Node is the the basis element of the ast
type Node interface {
	TokenLiteral() string
	String() string
}

// TokenLiteral impl for Program
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// **---TokenLiteral-implementations---** //

// TokenLiteral impl for LetStatement
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

// TokenLiteral impl for Identifier
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// TokenLiteral impl for ExpressionStatement
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

// TokenLiteral for ReturnStatement
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

// TokenLiteral for IntegerLiteral
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// TokenLiteral for FloatLiteral
func (f *FloatLiteral) TokenLiteral() string {
	return f.Token.Literal
}

// TokenLiteral for BooleanLiteral
func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}

// TokenLiteral for PrefixExpression
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

// TokenLiteral for InfixExpression
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

// **---String-implementations---** //

// String for Program
func (p *Program) String() string {
	var b bytes.Buffer

	for _, s := range p.Statements {
		b.WriteString(s.String())
	}

	return b.String()
}

// String for LetStatement
func (l *LetStatement) String() string {
	var b bytes.Buffer

	b.WriteString(l.TokenLiteral())
	b.WriteByte(' ')
	b.WriteString(l.Name.String())
	b.WriteString(" = ")

	if l.Value != nil {
		b.WriteString(l.Value.String())
	}

	b.WriteByte(';')

	return b.String()
}

// String for PrefixExpression
func (p *PrefixExpression) String() string {
	var b bytes.Buffer

	b.WriteByte('(')
	b.WriteString(p.Operator)
	b.WriteString(p.Right.String())
	b.WriteByte(')')

	return b.String()
}

// String for PrefixExpression
func (i *InfixExpression) String() string {
	var b bytes.Buffer

	b.WriteString("(")
	b.WriteString(i.Left.String())
	b.WriteByte(' ')
	b.WriteString(i.Operator)
	b.WriteByte(' ')
	b.WriteString(i.Right.String())
	b.WriteString(")")

	return b.String()
}

// String for ReturnStatement
func (r *ReturnStatement) String() string {
	var b bytes.Buffer

	b.WriteString(r.TokenLiteral())
	b.WriteByte(' ')

	if r.Value != nil {
		b.WriteString(r.Value.String())
	}

	b.WriteByte(';')

	return b.String()
}

func (i *Identifier) String() string {
	return i.Value
}

// String for ExpressionStatement
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

// String for IntegerLiteral
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

// String for FloatLiteral
func (f *FloatLiteral) String() string {
	return f.Token.Literal
}

func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

// Statement is the basis for a statment in the ast
type Statement interface {
	Node
	statementNode()
}

// **---statementNode-implementations---** //

func (l *LetStatement) statementNode()        {}
func (e *ExpressionStatement) statementNode() {}
func (r *ReturnStatement) statementNode()     {}

// Expression is the basis for a expression in the ast
type Expression interface {
	Node
	expressionNode()
}

// **---expressionNode-implementations---** //

func (i *Identifier) expressionNode()       {}
func (i *IntegerLiteral) expressionNode()   {}
func (f *FloatLiteral) expressionNode()     {}
func (p *PrefixExpression) expressionNode() {}
func (i *InfixExpression) expressionNode()  {}
func (b *BooleanLiteral) expressionNode()   {}
