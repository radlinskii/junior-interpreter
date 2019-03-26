package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/radlinskii/interpreter/evaluator"
	"github.com/radlinskii/interpreter/lexer"
	"github.com/radlinskii/interpreter/object"
	"github.com/radlinskii/interpreter/parser"
)

func main() {
	switch {
	case len(os.Args) == 1:
		fmt.Println("Please specify the file to be interpreted")
		os.Exit(1)
	case len(os.Args) > 2:
		fmt.Println("Please specify only one file to be interpreted")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	input := string(data)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) == 0 {
		env := object.NewEnvironment()

		evaluated := evaluator.EvalProgram(program, env)
		fmt.Println(evaluated)
	}
}
