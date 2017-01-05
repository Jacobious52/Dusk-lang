package check

import (
	"errors"
	"fmt"
	"jacob/black/token"
)

// Checker interface provides a standard function for performing a check in a token
// Implementor can keep state
// Tokens taken one at a time
type Checker interface {
	// Check single token
	Check(token.Token) error
	// Run at completion. Checker can do final checking
	Done() error
}

// Balanced checks if bracets, braces and stuff are Balanced
type Balanced struct {
	stack []token.Type
}

// Check function for Balanced
func (b *Balanced) Check(tok token.Token) error {

	switch tok.Type {
	case token.LBrace, token.LBracket, token.LParen:
		b.stack = append(b.stack, tok.Type)
	case token.RBrace:
		return b.popCheck(token.LBrace, tok)
	case token.RBracket:
		return b.popCheck(token.RBracket, tok)
	case token.RParen:
		return b.popCheck(token.LParen, tok)
	}

	return nil
}

func (b *Balanced) popCheck(match token.Type, tok token.Token) error {
	l := len(b.stack) - 1
	if l < 0 {
		return errors.New(fmt.Sprint("Extra ", tok.Literal))
	}
	if b.stack[l] != match {
		return errors.New(fmt.Sprint("Unbalanced ", token.LookupLiteral(b.stack[l]), ". Got ", tok.Literal))
	}
	b.stack = b.stack[:l]
	return nil
}

// Done function for Balanced
// Returns error if there are braces left in the stack
func (b *Balanced) Done() error {
	if len(b.stack) > 0 {
		return errors.New(fmt.Sprint("Unclosed ", token.LookupLiteral(b.stack[len(b.stack)-1])))
	}
	return nil
}
