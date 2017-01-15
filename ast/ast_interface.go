package ast

import (
	"bytes"
	"jacob/black/token"
	"strings"
)

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

// TokenLiteral for IfExpression
func (f *IfExpression) TokenLiteral() string {
	return f.Token.Literal
}

// TokenLiteral for FunctionLiteral
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

// TokenLiteral for BlockStatement
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
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

	b.WriteString("; ")

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

// String for IfExpression
func (f *IfExpression) String() string {
	var b bytes.Buffer

	b.WriteString("if ")
	b.WriteString(f.Cond.String())
	b.WriteByte(' ')
	b.WriteString(f.Do.String())

	if f.Else != nil {
		b.WriteString(" else ")
		b.WriteString(f.Else.String())
	}

	return b.String()
}

// String got FunctionLiteral
func (f *FunctionLiteral) String() string {
	var b bytes.Buffer

	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}

	b.WriteString(f.TokenLiteral())
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(f.TokenLiteral())
	b.WriteByte(' ')
	b.WriteString(f.Body.String())

	return b.String()
}

// String for BlockStatement
func (bs *BlockStatement) String() string {
	var b bytes.Buffer

	b.WriteString(bs.TokenLiteral())
	b.WriteByte(' ')

	for _, s := range bs.Statements {
		b.WriteString(s.String())
	}

	b.WriteByte(' ')

	if bs.Token.Type == token.LBrace {
		b.WriteByte('}')
	}

	return b.String()
}

func (c *CallExpression) String() string {
	var b bytes.Buffer

	args := []string{}
	for _, c := range c.Args {
		args = append(args, c.String())
	}

	b.WriteString(c.Func.String())
	b.WriteByte('(')
	b.WriteString(strings.Join(args, ", "))
	b.WriteByte(')')

	return b.String()
}

// String for Identifier
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
func (bs *BlockStatement) statementNode()     {}

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
func (f *IfExpression) expressionNode()     {}
func (f *FunctionLiteral) expressionNode()  {}
func (c *CallExpression) expressionNode()   {}
