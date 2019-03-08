package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/radlinskii/interpreter/lexer"
	"github.com/radlinskii/interpreter/token"
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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
