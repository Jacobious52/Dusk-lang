package main

import (
	"jacob/black/repl"
	"os"
)

func main() {
	repl.Run(os.Stdin, os.Stdout)
}
