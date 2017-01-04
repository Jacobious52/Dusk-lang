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

	curr int
	next int

	pos token.Position

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
func (l *Lexer) Init(filename string) {
	b, err := ioutil.ReadAll(l.src)
	if err != nil {
		panic(err)
	}
	l.buff = b

	l.pos.Filename = filename
	l.pos.Line = 1

	l.nextChar()
}

// NextToken returns the next token in the input stream
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.consumeWhitespace()

	switch l.char {
	case '=':
		tok = token.New(token.Assign, l.char, l.pos)
	case '+':
		tok = token.New(token.Plus, l.char, l.pos)
	case '-':
		tok = token.New(token.Minus, l.char, l.pos)
	case '*':
		tok = token.New(token.Times, l.char, l.pos)
	case '/':
		tok = token.New(token.Divide, l.char, l.pos)
	case '!':
		tok = token.New(token.Bang, l.char, l.pos)
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

func (l *Lexer) readIdentifier() token.Token {
	pos := l.pos

	p := l.curr
	for isLetter(l.char) {
		l.nextChar()
	}
	id := string(l.buff[p:l.curr])

	return token.Token{token.LookupIdenifier(id), id, pos}
}

// TODO: make parse double and int
func (l *Lexer) readNumber() token.Token {
	pos := l.pos

	p := l.curr
	for isDigit(l.char) {
		l.nextChar()
	}
	return token.Token{token.Int, string(l.buff[p:l.curr]), pos}
}

func (l *Lexer) nextChar() {
	if l.next >= len(l.buff) {
		l.char = 0
	} else {
		l.char = l.buff[l.next]
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
}
