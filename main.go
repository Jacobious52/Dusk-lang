package main

import (
	"jacob/black/repl"
	"os"
)

func main() {
	restart := true
	for restart {
		restart = repl.Run(os.Stdin, os.Stdout)
	}
}
