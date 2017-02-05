package compilejs

import (
	"fmt"
	"io"
	"jacob/dusk/lexer"
	"jacob/dusk/parser"
)

// Compile compiles to javascript
func Compile(in io.Reader, out io.Writer, name string) {
	l := lexer.WithReader(in, name)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(program.Statements) == 0 {
		return
	}

	if len(p.Errors()) != 0 {
		printErrors(out, p.Errors())
		return
	}

	fmt.Fprintln(out, program.String())
}

func printErrors(out io.Writer, errors []parser.Error) {
	for _, err := range errors {
		fmt.Fprintln(out, err.Pos, ":", err.Str)
	}
}
