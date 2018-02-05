package token

import "fmt"

// Type speficies a token type
type Type int

// Token is the struct that holds the stuff about tokens
type Token struct {
	Type
	Literal string
	Pos     Position
}

// New creates a new token
func New(t Type, literal byte, pos Position) Token {
	return Token{t, string(literal), pos}
}

// String for Token
func (t Token) String() string {
	lit := t.Type.String()
	if lit == "Unknown" {
		return t.Literal
	}
	return lit
}

// LookupIdenifier returns the Type for a Identifier string
func LookupIdenifier(id string) Type {
	if tok, ok := keywords[id]; ok {
		return tok
	}
	return Identifier
}

// Position is the location of a code point in the source
type Position struct {
	Filename string
	Offset   int
	Line     int
	Col      int
}

// Stringer method for Position
func (p Position) String() string {
	return fmt.Sprint(p.Filename, ":", p.Line, ":", p.Col)
}
