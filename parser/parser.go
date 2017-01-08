package parser

import (
	"jacob/black/ast"
	"jacob/black/lexer"
	"jacob/black/token"
)

// Parser parses into a ast from the lexer
type Parser struct {
	l *lexer.Lexer

	current token.Token
	next    token.Token

	errors []string
}

// New creates a new parser with the lexer l
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
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
	default:
		return nil
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

	// TODO: don't skip tokens
	for !p.expectNext(token.Terminator) {
		p.nextToken()
	}

	return let
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
	return false
}
