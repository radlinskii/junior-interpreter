# Junior Interpreter

[![GoDoc](https://godoc.org/github.com/radlinskii/interpreter?status.svg)](https://godoc.org/github.com/radlinskii/interpreter)
[![Build Status](https://travis-ci.com/radlinskii/interpreter.svg?branch=master)](https://travis-ci.com/radlinskii/interpreter)
[![Go Report Card](https://goreportcard.com/badge/github.com/radlinskii/interpreter)](https://goreportcard.com/report/github.com/radlinskii/interpreter)

Interpreter for programming language we named **Junior**.

## What is it

It's an interpreter written in Golang for programming language of our creation.<br />
It's a project for *Theory of Compilation* classes at the *AGH* university.<br />
It's based on [writing an interpreter in go](https://interpreterbook.com/) book.

You can test it in action [here](https://junior-interpreter-online.herokuapp.com/).

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
It's loosely typed but uses immutability. It has features like closures and IIFEs.
It's based on [Monkey programming language](https://interpreterbook.com/#the-monkey-programming-language).

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

### Grammar

#### Keywords

Reserved keywords of Junior:

`const, fun, return, if, else, true, false`

Reserved names of built-in functions:

`print, last, first, rest, len, push`

#### Builtins

Junior have some predefined functions that you can use.

1. `print` - prints given arguments to the output. returns null.
2. `len` - returns length of argument (array/string).
3. `first` - returns first element of an array.
4. `last` - returns last element of an array.
5. `rest` - returns all the elements of an array but the first one.
6. `push` - returns copy of given array with provided argument as the last element.

#### Literals

##### Booleans

Booleans are pretty straight forward.
Their values can only be either true or false.

```javascript
const truth = true;
const fact = truth != false; // true
```

##### Integers

Integers are as for now the only numeric values in Junior.
You can make every primitive mathematical operations on them.

```javascript
const number = 12;
const otherNumber = 34;

const sum = number + otherNumber; // 46
```

##### Strings

Strings are defined inside double-quotes.
As for now escaping double-quotes is not supported. But you don't need to escape e.g new lines.

```javascript
"The quick brown fox jumps over the lazy dog";
```

##### Functions

Functions in Junior are also treated as literals.
You can assign them to variables, store them in arrays or objects, pass them as arguments to other functions or immedietalt invoke them.

```javascript
const square = fun(x) {
    return x * x;
};

square(4); // 16
```

##### Arrays

Arrays in Junior are as immutable as any other literals.
They are not bound to one type but can store values of different types.
Arrays are indexed starting from 0.

```javascript
const arr = [true, 2, "three", fun(x) { return x * x; }];

arr[3](6); // 36
```

##### Hashes

Hashes are similar to Javascript's objects. Check them out:

```javascript
const obj = { name: "John Doe", age: "22", greet: fun(name) { return "Hi " + name + "! I'm John Doe"; } };

print(obj[name]); // John Doe
print(obj[greet]("Jane")); // Hi Jane! I'm John Doe 
```

#### Comments

Junior supports well known single and multi line comments.

```javascript
// This is a single line comment.
```

```javascript
/*
This is a
multi line comment.
 */
```

> Note: not terminated multi line comment will cause a parsing error.

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
