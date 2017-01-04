package lexer

import (
	"jacob/black/lexer"
	"jacob/black/token"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `let five = 5;
let ten = 10;

let add = |x, y| {
    x + y;
};

let result = add(five, ten);`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Identifier, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Terminator, ";"},
		{token.Let, "let"},
		{token.Identifier, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Terminator, ";"},
		{token.Let, "let"},
		{token.Identifier, "add"},
		{token.Assign, "="},
		{token.Bar, "|"},
		{token.Identifier, "x"},
		{token.Comma, ","},
		{token.Identifier, "y"},
		{token.Bar, "|"},
		{token.LBrace, "{"},
		{token.Identifier, "x"},
		{token.Plus, "+"},
		{token.Identifier, "y"},
		{token.Terminator, ";"},
		{token.RBrace, "}"},
		{token.Terminator, ";"},
		{token.Let, "let"},
		{token.Identifier, "result"},
		{token.Assign, "="},
		{token.Identifier, "add"},
		{token.LParen, "("},
		{token.Identifier, "five"},
		{token.Comma, ","},
		{token.Identifier, "ten"},
		{token.RParen, ")"},
		{token.Terminator, ";"},
		{token.EOF, string(0)},
	}

	l := lexer.WithString(input)
	l.Init()

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - Type wrong. expected=%q, got=%q",
				i, token.LookupLiteral(tt.expectedType), token.LookupLiteral(tok.Type))
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
