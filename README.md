# interpreter

[![GoDoc](https://godoc.org/github.com/radlinskii/interpreter?status.svg)](https://godoc.org/github.com/radlinskii/interpreter)
[![Build Status](https://travis-ci.com/radlinskii/interpreter.svg?branch=master)](https://travis-ci.com/radlinskii/interpreter)
[![Go Report Card](https://goreportcard.com/badge/github.com/radlinskii/interpreter)](https://goreportcard.com/report/github.com/radlinskii/interpreter)

Code of interpreter for C-based programming language following the [writing an interpreter in Go](https://interpreterbook.com/) book.

## TODO

- [ ] floats
- [x] counting rows
- [ ] counting columns
- [x] running with file as an argument
- [x] one-line comments
- [x] multiple-lines comments
- [x] return statement only permitted in function's body
- [x] return mandatory in function's body
- [x] return statement can be empty if we don't want to return anything
- [x] forbid reassigning constant
- [x] forbid redeclaring constant in one block
- [ ] ?? merge TRUE and FALSE tokens into one BOOLEAN token with different value
