package repl

import (
	"bufio"
	"fmt"
	"io"

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

		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg+"\n")
			}
			continue
		}

		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
