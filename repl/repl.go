package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"jacob/black/lexer"
	"jacob/black/token"
	"strings"
)

const (
	prompt = "=> "
	cont   = ".. "
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

	// Read until EOF
	for {
		fmt.Fprint(out, color(prompt, green))

		if ok := scanner.Scan(); !ok {
			return
		}

		// get current line
		line := scanner.Text()

		// if we open a function or map usinh {
		// continue reading into a buffer until we reach the maching }
		// then set this as the whole input
		if strings.HasSuffix(strings.TrimSpace(line), "{") {
			var indent int
			var b bytes.Buffer
			b.WriteString(line)

			for {
				fmt.Fprint(out, color(cont, green), strings.Repeat("\t", indent))

				if ok := scanner.Scan(); !ok {
					return
				}

				nextLine := scanner.Text()
				b.WriteString(nextLine)

				trimmed := strings.TrimSpace(nextLine)

				if strings.HasSuffix(trimmed, "}") {
					indent--
					if indent < 0 {
						break
					}
				} else if strings.HasSuffix(trimmed, "{") {
					// for if else case where end on close and start new
					if !strings.HasPrefix(trimmed, "}") {
						indent++
					}
				}
			}
			line = b.String()
		}

		l := lexer.WithString(line, "repl")

		for tok, err := l.Next(); tok.Type != token.EOF; tok, err = l.Next() {
			if err != nil {
				fmt.Fprintln(out, color("Error", red), color(tok.Pos, cyan), "-", err)
			}
			fmt.Fprintln(out, color("#", magneta), "\t", color(tok, yellow), "\t", color(tok.Type, cyan))
		}
	}
}
