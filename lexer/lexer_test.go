package lexer

import (
	"jacob/black/token"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `let five = 5
let ten = 10.2342

let add = |x, y| {
	x + y
}

let result = add(five, ten)

if result == 1
if result != 1
if !result < 3
let fail = 21
"foobar"
"foo bar"
"foo\tbar""foo\nbar"
`

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
		{token.Float, "10.2342"},
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
		{token.Let, "let"},
		{token.Identifier, "fail"},
		{token.Assign, "="},
		{token.Int, "21"},
		{token.Terminator, ";"},
		{token.String, "foobar"},
		{token.Terminator, ";"},
		{token.String, "foo bar"},
		{token.Terminator, ";"},
		{token.String, "foo\tbar"},
		{token.String, "foo\nbar"},
		{token.Terminator, ";"},
		{token.EOF, string(0)},
	}

	l := WithString(input, "lexer_test.go")

	for i, tt := range tests {
		tok, _ := l.Next()

		//t.Logf("%v\t\t%v\n", tok.Literal, tok.Pos)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - Type wrong. expected %q, got %q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong. expected %q, got %q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
