package repl

import (
	"bufio"
	"fmt"
	"io"
	"jacob/black/lexer"
	"jacob/black/token"
)

const (
	prompt = "> "
	intro  = "\033[2J\033[0;0HBlack Programming Langaguge (Repl). © Jacob Gonzalez 2017\n\n"
)

type colorCode int8

const (
	red = (iota + 31)
	green
	yellow
	blue
	magneta
	cyan
	white
)

func color(v interface{}, color colorCode) string {
	return fmt.Sprintf("\033[%vm%s\033[0m", color, v)
}

// Run starts the repl to read and run a line at a time
func Run(in io.Reader, out io.Writer) {

	fmt.Fprint(out, intro)

	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, color(prompt, green))

		if ok := scanner.Scan(); !ok {
			return
		}

		line := scanner.Text()

		l := lexer.WithString(line, "repl")

		for tok := l.Next(); tok.Type != token.EOF; tok = l.Next() {
			fmt.Fprintln(out, color("#", magneta), color(tok, yellow))
		}
	}
}
