package lexer

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"jacob/black/token"
	"strings"
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

	// stack for checking Balanced brackets, bracces..
	stack []token.Type
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

	l.last = token.Terminator

	l.nextChar()
}

// Next returns the next token in the input stream
func (l *Lexer) Next() (token.Token, error) {
	var tok token.Token
	var err error

	l.consumeWhitespace()

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			char := l.char
			l.nextChar()

			b := make([]byte, 2)
			b[0] = char
			b[1] = l.char

			tok = token.Token{Type: token.Equal, Literal: string(b), Pos: l.pos}
		} else {
			tok = token.New(token.Assign, l.char, l.pos)
		}
	case '+':
		if l.peekChar() == '=' {
			char := l.char
			l.nextChar()

			b := make([]byte, 2)
			b[0] = char
			b[1] = l.char

			tok = token.Token{Type: token.Inc, Literal: string(b), Pos: l.pos}
		} else {
			tok = token.New(token.Plus, l.char, l.pos)
		}
	case '-':
		if l.peekChar() == '=' {
			char := l.char
			l.nextChar()

			b := make([]byte, 2)
			b[0] = char
			b[1] = l.char

			tok = token.Token{Type: token.Dec, Literal: string(b), Pos: l.pos}
		} else {
			tok = token.New(token.Minus, l.char, l.pos)
		}
	case '*':
		tok = token.New(token.Times, l.char, l.pos)
	case '/':
		if l.peekChar() == '/' {
			l.consumeComment()
		} else {
			tok = token.New(token.Divide, l.char, l.pos)
		}
	case '^':
		tok = token.New(token.Exp, l.char, l.pos)
	case '%':
		tok = token.New(token.Mod, l.char, l.pos)
	case '!':
		if l.peekChar() == '=' {
			char := l.char
			l.nextChar()

			b := make([]byte, 2)
			b[0] = char
			b[1] = l.char

			tok = token.Token{Type: token.NotEqual, Literal: string(b), Pos: l.pos}
		} else {
			tok = token.New(token.Bang, l.char, l.pos)
		}
	case ':':
		tok = token.New(token.Continue, l.char, l.pos)
	case '<':
		tok = token.New(token.Less, l.char, l.pos)
	case '>':
		tok = token.New(token.Greater, l.char, l.pos)
	case '{':
		tok = token.New(token.LBrace, l.char, l.pos)
		l.stack = append(l.stack, tok.Type)
	case '}':
		tok = token.New(token.RBrace, l.char, l.pos)
		err = l.popCheck(token.LBrace, tok)
	case '(':
		tok = token.New(token.LParen, l.char, l.pos)
		l.stack = append(l.stack, tok.Type)
	case ')':
		tok = token.New(token.RParen, l.char, l.pos)
		err = l.popCheck(token.LParen, tok)
	case '[':
		tok = token.New(token.LBracket, l.char, l.pos)
		l.stack = append(l.stack, tok.Type)
		err = l.popCheck(token.LBracket, tok)
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
		if len(l.stack) > 0 {
			err = errors.New(fmt.Sprint("Unclosed ", l.stack[len(l.stack)-1]))
		}
	case '"':
		tok = l.readString()
	default:
		if isLetter(l.char) {
			return l.readIdentifier(), nil
		} else if isDigit(l.char) {
			return l.readNumber(), nil
		} else {
			tok = token.New(token.Illegal, l.char, l.pos)
		}
	}

	l.nextChar()

	l.last = tok.Type
	return tok, err
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
		// only add end of statment semi-colon if
		if l.char == '\n' && l.last != token.LBrace && l.last != token.Terminator {
			l.char = ';'
			return
		}
		l.nextChar()
	}
}

func (l *Lexer) consumeComment() {
	for l.char != '\n' {
		l.nextChar()
	}
}

func (l *Lexer) readIdentifier() token.Token {
	pos := l.pos
	p := l.curr

	// first must be _ or a-z
	l.nextChar()

	for isLetter(l.char) || isDigit(l.char) {
		l.nextChar()
	}
	id := string(l.buff[p:l.curr])

	l.last = token.Identifier
	return token.Token{Type: token.LookupIdenifier(id), Literal: id, Pos: pos}
}

// TODO: make parse double and int
func (l *Lexer) readNumber() token.Token {
	pos := l.pos

	p := l.curr
	for isDigit(l.char) {
		l.nextChar()
	}

	l.last = token.Int

	// read decimal number
	if l.char == '.' {
		l.nextChar()
		for isDigit(l.char) {
			l.nextChar()
		}
		l.last = token.Float
	}

	return token.Token{Type: l.last, Literal: string(l.buff[p:l.curr]), Pos: pos}
}

func (l *Lexer) readString() token.Token {
	pos := l.pos

	p := l.curr + 1
	for l.nextChar() != '"' {
	}

	str := string(l.buff[p:l.curr])
	r := strings.NewReplacer("\\t", "\t", "\\n", "\n")
	str = r.Replace(str)

	return token.Token{Type: token.String, Literal: str, Pos: pos}
}

func (l *Lexer) nextChar() byte {
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

	return l.char
}

func (l *Lexer) peekChar() byte {
	if l.next >= len(l.buff) {
		return 0
	}
	return l.buff[l.next]
}

// popCheck brace returning error if Unbalanced
func (l *Lexer) popCheck(match token.Type, tok token.Token) error {
	last := len(l.stack) - 1
	if last < 0 {
		return errors.New(fmt.Sprint("Extra ", tok.Literal))
	}
	if l.stack[last] != match {
		return errors.New(fmt.Sprint("Unbalanced ", l.stack[last], ". Got ", tok.Literal))
	}
	l.stack = l.stack[:last]
	return nil
}
