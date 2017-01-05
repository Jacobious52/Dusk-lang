package lexer

import (
	"io"
	"io/ioutil"
	"jacob/black/token"
)

// Lexer performs the tokenisation on a io.Reader
type Lexer struct {
	// src input buffer
	buff []byte

	// current index in buffer
	// is used to state the current place
	curr int
	// next index in buffer
	// is used to read the next character without affecting curr
	next int

	// position in the src visually
	pos token.Position

	// the current character the lexer is looking at
	char byte

	// the last token lexed
	// used for inserting semi-colon on line break
	last token.Type
}

// WithReader creates a new Lexer from the reader
func WithReader(reader io.Reader, filename string) *Lexer {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	l := &Lexer{buff: b}
	l.init(filename)

	return l
}

// WithString creates a new Lexer from a string
func WithString(str, filename string) *Lexer {
	l := &Lexer{buff: []byte(str)}
	l.init(filename)
	return l
}

// init sets the initial positons for the lexer
func (l *Lexer) init(filename string) {
	l.pos.Filename = filename
	l.pos.Line = 1

	l.nextChar()
}

// Next returns the next token in the input stream
func (l *Lexer) Next() token.Token {
	var tok token.Token

	l.consumeWhitespace()

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			char := l.char
			l.nextChar()

			b := make([]byte, 2)
			b[0] = char
			b[1] = l.char

			tok = token.Token{Type: token.Equal, Literal: string(b)}
		} else {
			tok = token.New(token.Assign, l.char, l.pos)
		}
	case '+':
		tok = token.New(token.Plus, l.char, l.pos)
	case '-':
		tok = token.New(token.Minus, l.char, l.pos)
	case '*':
		tok = token.New(token.Times, l.char, l.pos)
	case '/':
		tok = token.New(token.Divide, l.char, l.pos)
	case '!':
		if l.peekChar() == '=' {
			char := l.char
			l.nextChar()

			b := make([]byte, 2)
			b[0] = char
			b[1] = l.char

			tok = token.Token{Type: token.NotEqual, Literal: string(b)}
		} else {
			tok = token.New(token.Bang, l.char, l.pos)
		}
	case '<':
		tok = token.New(token.Less, l.char, l.pos)
	case '>':
		tok = token.New(token.Greater, l.char, l.pos)
	case '{':
		tok = token.New(token.LBrace, l.char, l.pos)
	case '}':
		tok = token.New(token.RBrace, l.char, l.pos)
	case '(':
		tok = token.New(token.LParen, l.char, l.pos)
	case ')':
		tok = token.New(token.RParen, l.char, l.pos)
	case '[':
		tok = token.New(token.LBracket, l.char, l.pos)
	case ']':
		tok = token.New(token.RBracket, l.char, l.pos)
	case '|':
		tok = token.New(token.Bar, l.char, l.pos)
	case ',':
		tok = token.New(token.Comma, l.char, l.pos)
	case ';':
		tok = token.New(token.Terminator, l.char, l.pos)
	case 0:
		tok = token.New(token.EOF, l.char, l.pos)
	default:
		if isLetter(l.char) {
			return l.readIdentifier()
		} else if isDigit(l.char) {
			return l.readNumber()
		} else {
			tok = token.New(token.Illegal, l.char, l.pos)
		}
	}

	l.nextChar()

	l.last = tok.Type
	return tok
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func (l *Lexer) consumeWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		// semi-colon insertion
		// only add eof of statment semi-colon if
		if l.char == '\n' && l.last != token.LBrace && l.last != token.Terminator {
			l.char = ';'
			return
		}
		l.nextChar()
	}
}

func (l *Lexer) readIdentifier() token.Token {
	pos := l.pos

	p := l.curr
	for isLetter(l.char) {
		l.nextChar()
	}
	id := string(l.buff[p:l.curr])

	l.last = token.Identifier
	return token.Token{token.LookupIdenifier(id), id, pos}
}

// TODO: make parse double and int
func (l *Lexer) readNumber() token.Token {
	pos := l.pos

	p := l.curr
	for isDigit(l.char) {
		l.nextChar()
	}

	l.last = token.Int
	return token.Token{token.Int, string(l.buff[p:l.curr]), pos}
}

func (l *Lexer) nextChar() {
	if l.next >= len(l.buff) {
		l.char = 0
	} else {
		l.char = l.buff[l.next]
	}

	l.curr = l.next
	l.next++

	// update position data
	l.pos.Col++
	l.pos.Offset = l.curr
	if l.char == '\n' {
		l.pos.Line++
		l.pos.Col = 0
	}

}

func (l *Lexer) peekChar() byte {
	if l.next >= len(l.buff) {
		return 0
	}
	return l.buff[l.next]
}
