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

let result = add(five, ten);

if result == 1;
if result != 1;
if !result < 3;`

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
		{token.If, "if"},
		{token.Identifier, "result"},
		{token.Equal, "=="},
		{token.Int, "1"},
		{token.Terminator, ";"},
		{token.If, "if"},
		{token.Identifier, "result"},
		{token.NotEqual, "!="},
		{token.Int, "1"},
		{token.Terminator, ";"},
		{token.If, "if"},
		{token.Bang, "!"},
		{token.Identifier, "result"},
		{token.Less, "<"},
		{token.Int, "3"},
		{token.Terminator, ";"},
		{token.EOF, string(0)},
	}

	l := lexer.WithString(input)
	l.Init("lexer_test.go")

	for i, tt := range tests {
		tok := l.NextToken()

		//t.Logf("%v\t\t%v\n", tok.Literal, tok.Pos)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - Type wrong. expected %q, got %q",
				i, token.LookupLiteral(tt.expectedType), token.LookupLiteral(tok.Type))
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong. expected %q, got %q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
