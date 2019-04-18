# Junior Interpreter

[![GoDoc](https://godoc.org/github.com/radlinskii/interpreter?status.svg)](https://godoc.org/github.com/radlinskii/interpreter)
[![Build Status](https://travis-ci.com/radlinskii/junior-interpreter.svg?branch=master)](https://travis-ci.com/radlinskii/junior-interpreter)
[![Go Report Card](https://goreportcard.com/badge/github.com/radlinskii/interpreter)](https://goreportcard.com/report/github.com/radlinskii/interpreter)
[![version](https://img.shields.io/github/release/radlinskii/junior-interpreter.svg)](https://img.shields.io/github/release/radlinskii/junior-interpreter.svg)

Interpreter for programming language named **Junior**.

## What is it

It's an interpreter written in Golang for programming language of our creation.<br />
It's a project for *Theory of Compilation* classes at the *AGH* university.<br />
It's based on [writing an interpreter in go](https://interpreterbook.com/) book.

You can test it in action in our online REPL [here](https://junior-interpreter-online.herokuapp.com/).

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

Junior is an imperative programming language. It derives from functional programming paradigm.
It's loosely typed but uses immutability. It has features like closures and IIFEs.
It's based on [Monkey programming language](https://interpreterbook.com/#the-monkey-programming-language).

### Table of contents
+ [Example program](#example-program)
+ [Keywords](#keywords)
+ [Statements](#statements)
  - [Const statement](#const-statement)
  - [Return statement](#return-statement)
  - [If statement](#if-statement)
  - [Expression Statement](#expression-statement)
+ [Expressions](#expressions)
  - [Literals](#literals)
    * [Booleans](#booleans)
    * [Integers](#integers)
    * [Strings](#strings)
    * [Functions](#functions)
    * [Arrays](#arrays)
    * [Hashes](#hashes)
  - [Operations](#operations)
    * [Logical](#logical)
    * [Mathematical](#mathematical-)
    * [Concatenation](#concatenation)
    * [Number Negation](#number-negation)
    * [Boolean Negation](#boolean-negation)
    * [Function Call](#function-call)
    * [Retrieving value with Index](#retrieving-value-with-index)
  - [Identifiers](#identifiers)
+ [Builtins](#builtins)
+ [Comments](#comments)
+ [Whitespaces](#whitespaces)
+ [Error handling](#error-handling)

### Example program

```javascript
const factorial = fun(x) {
    if (x < 1) {
        return 1;
    }
    return factorial(x - 1) * x;
};

factorial(5); // 120
```

You can find more examples in the `examples` directory.

### Keywords

Reserved keywords of Junior:

`const, fun, return, if, else, true, false`

Reserved names of built-in functions:

`print, last, first, rest, len, push`

### Statements

Junior program consists of instructions which are called statements.
In Junior statements are separated with semicolons.

#### Const statement

`const` `identifier` `=` `expression` `;`

Const statement binds the value evaluated from `expression` to a variable `identifier`.
Junior uses block scoping, there are three different kinds of scopes.
1. global scope
2. function scope
3. if/else scope

If variable is not found in the current scope the ancestor's scope is examined, if interpreter fails to find given identifier even in the global scope a semantic error is evaluated.
You cannot redeclare a variable that `identifier` represents in one scope.

#### Return statement

`return` `expression` `;` or `return` `;`

There are two rules when it comes to return statements in Junior:

1. Return statements are forbidden outside a function body.
2. Return statements are mandatory inside a function body.

> Note that you can omit an expression in return statement if you want your function to return `null`.

#### If statement

`if` `(` `condition` `)` `{` `consequence` `}`

or

`if` `(` `condition` `)` `{` `consequence` `}` `else` `{` `alternative` `}`

*If statement* evaluates statements in the *consequence* block if the *condition* was true.
If *condition* was false and *alternative* block is present it will get evaluated instead.

> Note in Junior `condition` must evaluate to a boolean, therefore this code:
` if (1) { print("1"); }` is not valid.

#### Expression Statement

In Junior every *expression* is also a *statement* therefore interpreter evaluates necessary expressions like e.g. function calls.

```javascript
1+1;
"Hello" + " World!";
print("Hello World!");
myAdder(40, 2);
```

> Note the above expressions are valid statements in Junior, but the first two don't make much sense outside the REPL, though.

### Expressions

In Junior, every operation and every literal is a valid expression and gets evaluated when interpreter is running.

#### Literals

In Junior every *literal* is an expression.
Only Booleans, Integers and Strings are "primitive" types that can be compared or treated as keys inside Hashes.

##### Booleans

Booleans are pretty straight forward.
Their values can only be either true or false.

```javascript
const truth = true;
const fact = truth != false; // true
```

##### Integers

Integers are as for now the only numeric values in Junior.
You can perform every primitive mathematical operations on them.

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

> Note: not terminating a string will cause a parsing error.

##### Functions

`fun` `(` `identifiers...` `)` `{` `statements...` `}`

Functions in Junior are also treated as literals. 
You can assign them to variables, store them in arrays or objects, pass them as arguments to other functions or immediately invoke them.
While evaluating function body, interpreter must encounter a return statement.
If your function is supposed to not return anything you can omit value in return statement.


```javascript
const square = fun(x) {
    return x * x;
};

square(4); // returns 16


const printSquare = fun(x) {
    print(x * x);
    
    return;
};

printSquare(4); // prints 16, returns null.
```

##### Arrays

`[` `expressions...` `]`

Arrays in Junior are as immutable as any other literals.
They are not bound to one type but can store values of different types.
Arrays are indexed starting from 0.

```javascript
const arr = [true, 2, "three", fun(x) { return x * x; }];

arr[3](6); // 36
```

##### Hashes

`{` `primitive type literal` `:` `expression` ... `}`

Hashes are maps with key: value pairs.
They are similar to Javascript's objects. Check them out:

```javascript
const obj = { "name": "John Doe", 4: "Been there.", "greet": fun(name) { return "Hi " + name + "! I'm John Doe"; } };

const greetStr = "greet";

print(obj[4]); // Been there.
print(obj[greetStr]("Jane")); // Hi Jane! I'm John Doe 
```

#### Operations

Junior supports many operations, from adding to numbers to retrieving value from a hash or array.
Here is a list of Junior's operations in order of their precedence.


##### Logical 

operators: `==`, `!=`, `>=`, `<=`, `>`, `<`

They evaluate and return logical value of expression they represent.
> Note that as for now they only support primitive types (booleans, integers, strings) as their operands.

##### Mathematical:

operators: `+`,`-`, `*`, `/`

Those operators return result of mathematical operation evaluated between their operands.
They only support integers as their operands.

```javascript
30 + 12;
84 / 2;
1 * 42;
42 - 0;
```

##### Concatenation

operator: `+`

Adds together two strings and returns the result.

```javascript
"Hello" + " World!";
```

##### Number Negation

operator: `-`

Prefixed operator for negating an integer.

```javascript
52 + -10;
```

##### Boolean Negation

operator: `!`

Prefixed operator for negating a boolean expression.

```javascript
const truth = true;
!truth;
```

##### Function Call

operators: `()`

To call a function put parenthesis after identifier or any other expression that evaluates to a function literal and pass arguments between them.

```javascript
const myWeirdFunction = fun(x) {
    return x + 2;
};

myWeirdFunction(40);

fun() {
    print("This is an IIFE!");
}();
```

##### Retrieving value with Index

operators: `[]`

The bracket operators are used to retrieve values from arrays and hashes.
They work just as in any other language.

```javascript
const myArray = [4, 2, 0];

myArray[0]; // 4

const theUniverse = { 42: "the answer", "isEarthFlat": false };

theUniverse[42];
theUniverse["isEarthFlat"];
```

#### Identifiers

Identifiers are also treated as expressions.
They evaluate to expression they are bound to.
You can't redeclare a variable inside it's scope but you can overwrite, an identifier that was declared inside scope of one of an ancestors of current scope.
```javascript
const randomNumber = 40;
const two = 2;

randomNumber + two;
```

### Builtins

Junior have some predefined functions that you can use.

1. `print(values...)` - prints given arguments to the output, returns null.
2. `len(array|string)` - returns length of argument (array or string).
3. `first(array)` - returns first element of an array.
4. `last(array)` - returns last element of given array.
5. `rest(array)` - returns all the elements of given array but the first one.
6. `push(array|value)` - returns copy of given array with provided argument as the last element.


### Comments

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


### Whitespaces

From the interpreter's perspective *whitespaces* are meaningless, but you should always focus on your code's readability.


### Error handling

1. Every **Lexical error**, e.g. *invalid token*, stops interpreter from parsing the program.
2. **Syntax errors**, e.g. *missing semicolon*, are collected through parsing and printed after parsing process is finished. They prevent program from being evaluated.
3. Any **Semantic error**, e.g. *type incompatibility*, or **Evaluation errors**, e.g. *index out of boundaries*, stops evaluation of the program.

## Installation and development

1. `clone` the *repository*
2. `go get` the *dependencies*
3. `export` the *PORT* environment variable
4. `go run` the *main.go* file

## Contributing

Found a bug or typo? Create an issue [here](https://github.com/radlinskii/junior-interpreter/issues/new).
You are also welcome to fork this repository and create PRs, but remember to create an issue and assign it to yourself first.

## Authors

[@radlinskii](https://github.com/radlinskii)
[@agnieszka-miszkurka](https://github.com/agnieszka-miszkurka)
