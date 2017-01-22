package parser

import (
	"fmt"
	"jacob/dusk/ast"
	"jacob/dusk/lexer"
	"jacob/dusk/token"
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
	assign                // =
	sum                   // + -
	product               // * /
	exp                   // ^ %
	prefix                // -X or !X
	call                  // f(x)
	index                 // a[i]
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
	token.LParen:   call,
	token.Bang:     call,
	token.Assign:   assign,
	token.Exp:      exp,
	token.Mod:      exp,
	token.Inc:      assign,
	token.Dec:      assign,
	token.LBracket: index,
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
	p.registerPrefix(token.If, p.parseIfExpression)
	p.registerPrefix(token.Bar, p.parseFunctionLiteral)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.LBracket, p.parseArrayLiteral)

	p.infixParseFn = make(map[token.Type]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Divide, p.parseInfixExpression)
	p.registerInfix(token.Times, p.parseInfixExpression)
	p.registerInfix(token.Exp, p.parseInfixExpression)
	p.registerInfix(token.Mod, p.parseInfixExpression)
	p.registerInfix(token.Equal, p.parseInfixExpression)
	p.registerInfix(token.Inc, p.parseInfixExpression)
	p.registerInfix(token.Dec, p.parseInfixExpression)
	p.registerInfix(token.NotEqual, p.parseInfixExpression)
	p.registerInfix(token.Assign, p.parseInfixExpression)
	p.registerInfix(token.Less, p.parseInfixExpression)
	p.registerInfix(token.Greater, p.parseInfixExpression)
	p.registerInfix(token.LParen, p.parseCallExpression)
	p.registerInfix(token.Bang, p.parseCallExpression)
	p.registerInfix(token.LBracket, p.parseIndexExpression)

	// load the next token and current token
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
	msg := fmt.Sprintf("expected next token to be '%s', got '%s' instead", t, p.next)
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

	let.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

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
		p.newError(fmt.Sprintf("'%s' is not a valid operator", p.current))
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
	expr := &ast.PrefixExpression{Token: p.current, Operator: p.current.Type}

	p.nextToken()
	expr.Right = p.parseExpression(prefix)

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// += and -= operators patch ast in real time
	// expand a += 3 to a = a + 3
	if p.current.Type == token.Inc || p.current.Type == token.Dec {

		// top level is a = (...)
		expr := &ast.InfixExpression{Token: token.Token{Type: token.Assign, Literal: token.Assign.String(), Pos: p.current.Pos}, Left: left, Operator: token.Assign}

		expanded := token.Plus
		if p.current.Type == token.Dec {
			expanded = token.Minus
		}

		// next level is ... = a + expr..
		exprPlus := &ast.InfixExpression{Token: token.Token{Type: expanded, Literal: expanded.String(), Pos: p.current.Pos}, Left: left, Operator: expanded}
		expr.Right = exprPlus

		prec := p.currentPrecedence()
		p.nextToken()
		exprPlus.Right = p.parseExpression(prec)

		return expr
	}

	expr := &ast.InfixExpression{Token: p.current, Left: left, Operator: p.current.Type}

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

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.current}

	p.nextToken()
	expr.Cond = p.parseExpression(lowest)

	// check if with mult statement or single statement
	if !(p.nextIs(token.LBrace) || p.nextIs(token.Continue)) {
		p.newError(fmt.Sprintf("expected '{' or ':' following if statement, got '%s' instead", p.next))
		return nil
	}

	// goto the { or : and begin the block statment
	p.nextToken()
	expr.Do = p.parseBlockStatement()

	if p.nextIs(token.Else) {
		p.nextToken()

		// current is else.
		// optional { : or none
		if p.nextIs(token.LBrace) || p.nextIs(token.Continue) {
			p.nextToken()
		}

		// goto { or none and parse block
		expr.Else = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	f := &ast.FunctionLiteral{Token: p.current}

	// current is | or !
	f.Params = p.parseFunctionParams()

	// current is ending |
	// optional {
	if p.nextIs(token.LBrace) || p.nextIs(token.Continue) {
		p.nextToken()
	}

	// if current is not { function will be single statment
	f.Body = p.parseBlockStatement()

	return f
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	ids := []*ast.Identifier{}

	// capture empty args ! and just return
	if p.currentIs(token.Bang) {
		return ids
	}

	// '||' empty params
	if p.nextIs(token.Bar) {
		p.nextToken()
		return ids
	}

	// that means at least one param. get it
	if !p.expectNext(token.Identifier) {
		return nil
	}
	ids = append(ids, &ast.Identifier{Token: p.current, Value: p.current.Literal})

	// keep getting params until no more commas
	for p.nextIs(token.Comma) {
		// swollow comma
		p.nextToken()

		// param must be id
		if !p.expectNext(token.Identifier) {
			return nil
		}
		ids = append(ids, &ast.Identifier{Token: p.current, Value: p.current.Literal})
	}

	// must end with bar
	if !p.expectNext(token.Bar) {
		return nil
	}

	return ids
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// keep the leading token to tell us if -> or {
	leading := p.current
	p.nextToken()

	block := &ast.BlockStatement{Token: leading}
	block.Statements = []ast.Statement{}

	// catch empty statement
	if p.currentIs(token.RBrace) || p.currentIs(token.Terminator) {
		return block
	}

	// try parse the first statement. always should be one statment for ->
	// don't go next token because } might or might not exist
	s := p.parseStatement()
	if s != nil {
		block.Statements = append(block.Statements, s)
	}

	// if { then keep adding statemnts until }
	if leading.Type == token.LBrace {
		// means last statement ended on }. skip it
		p.nextToken()
		// keep getting statemnts until we reach final }
		for !p.currentIs(token.RBrace) {
			s := p.parseStatement()
			if s != nil {
				block.Statements = append(block.Statements, s)
			}
			// skip the }
			p.nextToken()
		}
	}

	return block
}

func (p *Parser) parseCallExpression(f ast.Expression) ast.Expression {
	call := &ast.CallExpression{Token: p.current, Func: f}
	call.Args = p.parseCallArgs()
	return call
}

func (p *Parser) parseCallArgs() []ast.Expression {
	args := []ast.Expression{}

	// capture empty args ! and just return
	if p.currentIs(token.Bang) {
		return args
	}

	// '()' empty params
	if p.nextIs(token.RParen) {
		p.nextToken()
		return args
	}
	p.nextToken()

	// that means at least one arg. get it
	args = append(args, p.parseExpression(lowest))

	// keep getting params until no more commas
	for p.nextIs(token.Comma) {
		// swollow comma
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(lowest))
	}

	// must end with )
	if !p.expectNext(token.RParen) {
		return nil
	}

	return args
}

func (p *Parser) parseIdentifier() ast.Expression {

	if p.nextIs(token.Dot) {
		id := &ast.AccessIdentifier{Token: p.current}
		values := []string{p.current.Literal}

		for p.nextIs(token.Dot) {
			p.nextToken()
			p.nextToken()
			values = append(values, p.current.Literal)
		}

		id.Values = values
		return id
	}

	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.current}

	if val, err := strconv.ParseInt(p.current.Literal, 0, 64); err == nil {
		lit.Value = val
		return lit
	}

	msg := fmt.Sprintf("could not parse '%s' as Integer", p.current.Literal)
	p.newError(msg)
	return nil
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.current}

	if val, err := strconv.ParseFloat(p.current.Literal, 64); err == nil {
		lit.Value = val
		return lit
	}

	msg := fmt.Sprintf("could not parse '%s' as Float", p.current.Literal)
	p.newError(msg)
	return nil
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.current}
	arr.Elements = p.parseListElems()
	return arr
}

func (p *Parser) parseListElems() []ast.Expression {
	args := []ast.Expression{}

	// '[]' empty params
	if p.nextIs(token.RBracket) {
		p.nextToken()
		return args
	}
	p.nextToken()

	// that means at least one arg. get it
	args = append(args, p.parseExpression(lowest))

	// keep getting params until no more commas
	for p.nextIs(token.Comma) {
		// swollow comma
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(lowest))
	}

	// must end with ]
	if !p.expectNext(token.RBracket) {
		return nil
	}

	return args
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Token: p.current, Left: left}
	p.nextToken()

	expr.Index = p.parseExpression(lowest)

	if !p.expectNext(token.RBracket) {
		return nil
	}

	return expr
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.BooleanLiteral{Token: p.current, Value: p.currentIs(token.True)}
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
