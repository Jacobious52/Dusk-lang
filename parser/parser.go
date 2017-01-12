package parser

import (
	"fmt"
	"jacob/black/ast"
	"jacob/black/lexer"
	"jacob/black/token"
	"strconv"
)

// Error holds a parser Error
// the positon it happended
// a message
type Error struct {
	Str string
	Pos token.Position
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type precedence int

const (
	lowest     precedence = (iota + 1)
	equals                // == !=
	inequality            // < >
	sum                   // + -
	product               // * /
	prefix                // -X or !X
	call                  // f(x)
)

var precedences = map[token.Type]precedence{
	token.Plus:     sum,
	token.Minus:    sum,
	token.Equal:    equals,
	token.NotEqual: equals,
	token.Less:     inequality,
	token.Greater:  inequality,
	token.Divide:   product,
	token.Times:    product,
}

// Parser parses into a ast from the lexer
type Parser struct {
	l *lexer.Lexer

	current token.Token
	next    token.Token

	errors []Error

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFn   map[token.Type]infixParseFn
}

// New creates a new parser with the lexer l
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.Identifier, p.parseIdentifier)
	p.registerPrefix(token.Int, p.parseIntegerLiteral)
	p.registerPrefix(token.Float, p.parseFloatLiteral)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBooleanExpression)
	p.registerPrefix(token.False, p.parseBooleanExpression)
	p.registerPrefix(token.LParen, p.parseGroupedExpression)

	p.infixParseFn = make(map[token.Type]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Divide, p.parseInfixExpression)
	p.registerInfix(token.Times, p.parseInfixExpression)
	p.registerInfix(token.Equal, p.parseInfixExpression)
	p.registerInfix(token.NotEqual, p.parseInfixExpression)
	p.registerInfix(token.Less, p.parseInfixExpression)
	p.registerInfix(token.Greater, p.parseInfixExpression)

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

func (p *Parser) currentPrecedence() precedence {
	if p, ok := precedences[p.current.Type]; ok {
		return p
	}

	return lowest
}

func (p *Parser) nextPrecedence() precedence {
	if p, ok := precedences[p.next.Type]; ok {
		return p
	}

	return lowest
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
	p.errors = append(p.errors, Error{str, p.current.Pos})
}

func (p *Parser) newPeekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.next.Type)
	p.errors = append(p.errors, Error{msg, p.current.Pos})
}

// Errors returns all the errors the parser encountered
func (p *Parser) Errors() []Error {
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

	p.nextToken()

	let.Value = p.parseExpression(lowest)

	if p.nextIs(token.Terminator) {
		p.nextToken()
	}

	return let
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	ret := &ast.ReturnStatement{Token: p.current}

	p.nextToken()

	ret.Value = p.parseExpression(lowest)

	if p.nextIs(token.Terminator) {
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
	// try parse prefix expression first
	prefixParser, ok := p.prefixParseFns[p.current.Type]
	if !ok {
		p.newError(fmt.Sprintf("no operand for prefix operator '%s' found", p.current.Type))
		return nil
	}
	leftExpr := prefixParser()

	// parse infix expression unitl reach a higher precedence
	for prec < p.nextPrecedence() {
		// break if next token is not an infix operator
		// this includes semi colon
		infixParser, ok := p.infixParseFn[p.next.Type]
		if !ok {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infixParser(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{Token: p.current, Operator: p.current.Literal}
	p.nextToken()
	expr.Right = p.parseExpression(prefix)
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{Token: p.current, Left: left, Operator: p.current.Literal}

	prec := p.currentPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(prec)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	expr := p.parseExpression(lowest)
	if !p.expectNext(token.RParen) {
		return nil
	}
	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.current}

	if val, err := strconv.ParseInt(p.current.Literal, 0, 64); err == nil {
		lit.Value = val
		return lit
	}

	msg := fmt.Sprintf("could not parse %q as Integer", p.current.Literal)
	p.newError(msg)

	return nil
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.current}

	if val, err := strconv.ParseFloat(p.current.Literal, 64); err == nil {
		lit.Value = val
		return lit
	}

	msg := fmt.Sprintf("could not parse %s as Float", p.current.Literal)
	p.newError(msg)

	return nil
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.BooleanLiteral{p.current, p.currentIs(token.True)}
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
