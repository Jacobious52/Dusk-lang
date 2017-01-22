package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"jacob/dusk/eval"
	"jacob/dusk/lexer"
	"jacob/dusk/object"
	"jacob/dusk/parser"
	"strings"
)

const (
	prompt = "| "
	intro  = "\033[2J\033[0;0HDusk Programming Langaguge (Repl). © Jacob Gonzalez 2017\n\n"
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
func Run(in io.Reader, out io.Writer) bool {

	fmt.Fprint(out, intro)

	scanner := bufio.NewScanner(in)

	env := object.NewEnvironment()

	// Read until EOF
	for {
		lineNum := 1
		fmt.Fprint(out, lineNum, color(prompt, green))

		if ok := scanner.Scan(); !ok {
			return false
		}

		// get current line
		line := scanner.Text()

		switch line {
		case ":r":
			return true
		case ":x", ":q", ":e":
			return false
		case ":c":
			fmt.Fprint(out, intro)
			continue
		}

		var b bytes.Buffer
		b.WriteString(line)
		b.WriteByte('\n')

		// if we open a function or map usinh {
		// continue reading into a buffer until we reach the maching }
		// then set this as the whole input
		if strings.HasSuffix(strings.TrimSpace(line), "{") {
			var indent int
			for {
				lineNum++
				fmt.Fprint(out, lineNum, color(prompt, blue), strings.Repeat("\t", indent))

				if ok := scanner.Scan(); !ok {
					return false
				}

				nextLine := scanner.Text()
				b.WriteString(nextLine)
				b.WriteByte('\n')

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
		}

		l := lexer.WithString(b.String(), "repl")
		p := parser.New(l)

		program := p.ParseProgram()

		if len(program.Statements) == 0 {
			continue
		}

		if len(p.Errors()) != 0 {
			printErrors(out, p.Errors())
			continue
		}

		result := eval.Eval(program, env)

		if result != nil && result.Type() != object.NilType {
			fmt.Fprintln(out, "", color(prompt, magneta), "\t", color(strings.Replace(result.String(), "\n", fmt.Sprint("\n ", color(prompt, magneta), " \t "), -1), yellow))
		}
		fmt.Fprint(out, "\n")
	}
}

func printErrors(out io.Writer, errors []parser.Error) {
	for _, err := range errors {
		fmt.Fprintln(out, "", color(prompt, red), "\t", color(fmt.Sprint(err.Pos, ":"), red), err.Str)
	}
	fmt.Fprint(out, "\n")
}
