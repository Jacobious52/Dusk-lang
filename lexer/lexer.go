package lexer

import (
	"io"
	"io/ioutil"
	"jacob/black/token"
	"strings"
)

// Lexer performs the tokenisation on a io.Reader
type Lexer struct {
	src  io.Reader
	buff []byte

	pos  int
	next int
	line int

	char byte
}

// WithReader creates a new Lexer from the reader
func WithReader(reader io.Reader) *Lexer {
	l := &Lexer{src: reader}
	return l
}

// WithString creates a new Lexer from a string
func WithString(str string) *Lexer {
	return WithReader(strings.NewReader(str))
}

// Init readers the reader src into a buffer
func (l *Lexer) Init() {
	b, err := ioutil.ReadAll(l.src)
	if err != nil {
		panic(err)
	}
	l.buff = b

	l.nextChar()
}

// NextToken returns the next token in the input stream
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.consumeWhitespace()

	switch l.char {
	case '=':
		tok = token.New(token.Assign, l.char)
	case '+':
		tok = token.New(token.Plus, l.char)
	case '-':
		tok = token.New(token.Minus, l.char)
	case '*':
		tok = token.New(token.Times, l.char)
	case '/':
		tok = token.New(token.Divide, l.char)
	case '!':
		tok = token.New(token.Bang, l.char)
	case '<':
		tok = token.New(token.Less, l.char)
	case '>':
		tok = token.New(token.Greater, l.char)
	case '{':
		tok = token.New(token.LBrace, l.char)
	case '}':
		tok = token.New(token.RBrace, l.char)
	case '(':
		tok = token.New(token.LParen, l.char)
	case ')':
		tok = token.New(token.RParen, l.char)
	case '[':
		tok = token.New(token.LBracket, l.char)
	case ']':
		tok = token.New(token.RBracket, l.char)
	case '|':
		tok = token.New(token.Bar, l.char)
	case ',':
		tok = token.New(token.Comma, l.char)
	case ';':
		tok = token.New(token.Terminator, l.char)
	case 0:
		tok = token.New(token.EOF, l.char)
	default:
		if isLetter(l.char) {
			tok.Literal, tok.Type = l.readIdentifier()
			return tok
		} else if isDigit(l.char) {
			tok.Literal, tok.Type = l.readNumber()
			return tok
		} else {
			tok = token.New(token.Illegal, l.char)
		}
	}

	l.nextChar()

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
		l.nextChar()
	}
}

func (l *Lexer) readIdentifier() (string, token.Type) {
	p := l.pos
	for isLetter(l.char) {
		l.nextChar()
	}
	id := string(l.buff[p:l.pos])
	return id, token.LookupIdenifier(id)
}

// TODO: make parse double and int
func (l *Lexer) readNumber() (string, token.Type) {
	p := l.pos
	for isDigit(l.char) {
		l.nextChar()
	}
	return string(l.buff[p:l.pos]), token.Int
}

func (l *Lexer) nextChar() {
	if l.next >= len(l.buff) {
		l.char = 0
	} else {
		l.char = l.buff[l.next]
		l.pos = l.next
		l.next++
	}
}
