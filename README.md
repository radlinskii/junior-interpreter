# Junior Interpreter

[![GoDoc](https://godoc.org/github.com/radlinskii/interpreter?status.svg)](https://godoc.org/github.com/radlinskii/interpreter)
[![Build Status](https://travis-ci.com/radlinskii/interpreter.svg?branch=master)](https://travis-ci.com/radlinskii/interpreter)
[![Go Report Card](https://goreportcard.com/badge/github.com/radlinskii/interpreter)](https://goreportcard.com/report/github.com/radlinskii/interpreter)

Interpreter for programming language we named **Junior**.

## What is it

It's an interpreter written in Golang for programming language of our creation.<br />
It's a project for *Theory of Compilation* classes at the *AGH* university.<br />
It's based on [writing an interpreter in go](https://interpreterbook.com/) book.

You can use and test it in action [here](https://junior-interpreter-online.herokuapp.com/).

## How it Works

It divides interpreting the Junior's code into 3 parts.

1. Lexer - performing lexical analysis 
2. Parser - syntax analysis and building Abstract Syntax Tree.
3. Evaluator - traversing AST and evaluating the program.

|  | input | output |
| :---: | :---: | :---: |
| **Lexer** | Junior's program code | tokens |
| **Parser** | tokens | Abstract Syntax Tree nodes |
| **Evaluator** | Abstract Syntax Tree nodes | evaluated program statements |

## Junior Language Specification

Junior is an imperative programming language. It derive from functional programming paradigm.
It's loosely typed but uses immutability. It has features like closures and IIFEs. It's based on [Monkey programming language](https://interpreterbook.com/#the-monkey-programming-language).

Example program:
```javascript
const makeConcat = fun(y) {
    return fun(x) {
        return x + y;
    };
};

const addWorld = makeConcat(" World!");

addWorld("Hello"); // Hello World!
```

You can find more examples in the `examples` directory.

# TODO*

### Types

int. string, etc. functions ??? you can assign them to a variable so probably...
arrays, hashes...

### Expressions

things you can bind to a variable
literals : "aa";
expressions: 1 + 3
function calls: const a = first([1,3,5]);

functions can do it!

### Statements

things you can't assign to a constant:

1. const foo = "bar";

2. if (true) {
    print("foo");
}

3. return "bar";


### Builtins

first,last etc.

### features

closure
IIFE

# END OF TODO*

### Error handling

1. Every **Lexical error**, e.g. *invalid token*, stops program parsing.
2. **Syntax errors**, e.g. *missing semicolon*, are collected through parsing and printed after parsing process is finished. They prevent program from being evaluated.
3. Any **Semantic error**, e.g. *type incompatibility*, or **Evaluation errors**, e.g. *index out of boundaries*, stops evaluation of the program.

## Installation and development

1. `clone` the *repository*
2. `go get` the *dependencies*
3. `export` the *PORT* environment variable
4. `go run` the *main.go* file

## Contributing

Feel free to post issues.
You are also welcome fork the repo and create PRs, but remember to create an issue and assign it to yourself first.

## Authors

[@radlinskii](https://github.com/radlinskii)
[@agnieszka-miszkurka](https://github.com/agnieszka-miszkurka)
