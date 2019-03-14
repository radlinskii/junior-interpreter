package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/radlinskii/interpreter/evaluator"
	"github.com/radlinskii/interpreter/lexer"
	"github.com/radlinskii/interpreter/object"
	"github.com/radlinskii/interpreter/parser"
)

func main() {
	data, err := ioutil.ReadFile("example.monkey")
	if err != nil {
		log.Fatalf(err.Error())
	}
	input := string(data)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Println("ERROR: " + msg + "\n")
		}
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}
