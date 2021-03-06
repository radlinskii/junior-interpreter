package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/radlinskii/interpreter/object"

	"github.com/radlinskii/interpreter/evaluator"

	"github.com/radlinskii/interpreter/lexer"
	"github.com/radlinskii/interpreter/parser"
)

// PROMPT defines how the REPL's prompt will look like.
const PROMPT = "👉  "

// Start runs the REPL loop.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) == 0 {
			evaluated := evaluator.EvalProgram(program, env)
			fmt.Println(evaluated)
		}

	}
}

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)

	fmt.Println("Feel free to type in commands")
	Start(os.Stdin, os.Stdout)
}
