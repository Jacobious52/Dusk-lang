package main

import (
	"jacob/dusk/repl"
	"os"
)

func main() {
	restart := true
	for restart {
		restart = repl.Run(os.Stdin, os.Stdout)
	}
}
