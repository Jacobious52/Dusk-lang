package main

import (
	"fmt"
	"jacob/dusk/compile-js"
	"jacob/dusk/repl"
	"jacob/dusk/run"
	"os"
)

func main() {
	switch len(os.Args) {
	case 1: // run repl
		restart := true
		for restart {
			restart = repl.Run(os.Stdin, os.Stdout)
		}
	case 2: // run file
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		defer file.Close()

		run.Run(file, os.Stdout, os.Args[1])
	case 3: // run file with js
		if os.Args[1] == "--js" {
			file, err := os.Open(os.Args[2])
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
			defer file.Close()

			compilejs.Compile(file, os.Stdout, os.Args[2])
		}
	}
}
