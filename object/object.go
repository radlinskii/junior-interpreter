package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/radlinskii/interpreter/ast"
)

// Type represents different object types.
type Type string

const (
	// INTEGER object type
	INTEGER = "INTEGER"
	// BOOLEAN object type
	BOOLEAN = "BOOLEAN"
	// STRING object type
	STRING = "STRING"
	// NULL object type
	NULL = "NULL"
	// RETURN object wrapper type
	RETURN = "RETURN"
	// ERROR object type
	ERROR = "ERROR"
	// FUNCTION object type
	FUNCTION = "FUNCTION"
)

// Object interface is implemented by the objects.
type Object interface {
	Inspect() string
	Type() Type
}

// Integer object.
type Integer struct {
	Value int64
}

// Inspect returns value of an integer.
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type returns the integer type.
func (i *Integer) Type() Type {
	return INTEGER
}

// Boolean object.
type Boolean struct {
	Value bool
}

// Inspect returns value of a boolean.
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Type returns the boolean type.
func (b *Boolean) Type() Type {
	return BOOLEAN
}

// String object.
type String struct {
	Value string
}

// Inspect returns value of a string.
func (s *String) Inspect() string {
	return s.Value
}

// Type returns the string type.
func (s *String) Type() Type {
	return STRING
}

// Null object.
type Null struct{}

// Inspect returns null.
func (n *Null) Inspect() string {
	return "null"
}

// Type returns the null object type.
func (n *Null) Type() Type {
	return NULL
}

// Return object is a wrapper to a object that gets returned.
type Return struct {
	Value Object
}

// Inspect returns the value of object to be returned.
func (rv *Return) Inspect() string {
	return rv.Value.Inspect()
}

// Type returns the Return object type.
func (rv *Return) Type() Type {
	return RETURN
}

// Error object.
type Error struct {
	Message string
}

// Inspect returns error message.
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// Type returns the Error object type.
func (e *Error) Type() Type {
	return ERROR
}

// Function object.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Inspect returns the Function object image.
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fun(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

// Type returns the Function object type.
func (f *Function) Type() Type {
	return FUNCTION
}

// Environment is a map of known objects.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment returns new Environment instance
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment returns new Environment instance
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get returns value of given key from Enviroment's map.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set puts the value of given key in Enviroment's map.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
