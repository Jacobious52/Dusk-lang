package parser

import (
	"fmt"
	"jacob/black/ast"
	"jacob/black/lexer"
	"jacob/black/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type precedence int

const (
	lowest     precedence = (iota + 1)
	equals                // ==
	inequality            // < >
	sum                   // +
	product               // *
	prefix                // -X or !X
	call                  // f(x)
)

// Parser parses into a ast from the lexer
type Parser struct {
	l *lexer.Lexer

	current token.Token
	next    token.Token

	errors []string

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFn   map[token.Type]infixParseFn
}

// New creates a new parser with the lexer l
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.Identifier, p.parseIdentifier)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(t token.Type, f prefixParseFn) {
	p.prefixParseFns[t] = f
}

func (p *Parser) registerInfix(t token.Type, f infixParseFn) {
	p.infixParseFn[t] = f
}

func (p *Parser) nextToken() {
	var err error

	p.current = p.next
	p.next, err = p.l.Next()
	if err != nil {
		p.newError(err.Error())
	}
}

func (p *Parser) newError(str string) {
	p.errors = append(p.errors, str)
}

func (p *Parser) newPeekError(t token.Type) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.next.Type)
	p.errors = append(p.errors, msg)
}

// Errors returns all the errors the parser encountered
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram parses a whole program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.current.Type != token.EOF {
		if statement := p.parseStatement(); statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	let := &ast.LetStatement{Token: p.current}

	if !p.expectNext(token.Identifier) {
		return nil
	}

	let.Name = &ast.Identifier{p.current, p.current.Literal}

	if !p.expectNext(token.Assign) {
		return nil
	}

	//let.Value = p.parseExpressionStatement().Expression

	// TODO: don't skip tokens
	for !p.currentIs(token.Terminator) {
		p.nextToken()
	}

	return let
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	ret := &ast.ReturnStatement{Token: p.current}

	p.nextToken()

	//ret.Value = p.parseExpressionStatement().Expression

	// TODO: don't skip tokens
	for !p.currentIs(token.Terminator) {
		p.nextToken()
	}

	return ret
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expr := &ast.ExpressionStatement{Token: p.current}

	expr.Expression = p.parseExpression(lowest)

	if p.nextIs(token.Terminator) {
		p.nextToken()
	}

	return expr
}

func (p *Parser) parseExpression(prec precedence) ast.Expression {
	if prefixParser, ok := p.prefixParseFns[p.current.Type]; ok {
		leftExpr := prefixParser()
		return leftExpr
	}
	return nil
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) currentIs(t token.Type) bool {
	return p.current.Type == t
}

func (p *Parser) nextIs(t token.Type) bool {
	return p.next.Type == t
}

func (p *Parser) expectNext(t token.Type) bool {
	if p.nextIs(t) {
		p.nextToken()
		return true
	}
	p.newPeekError(t)
	return false
}
