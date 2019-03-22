package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
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
	// VOID object type
	VOID = "VOID"
	// RETURN object wrapper type
	RETURN = "RETURN"
	// ERROR object type
	ERROR = "ERROR"
	// FUNCTION object type
	FUNCTION = "FUNCTION"
	// BUILTIN object type
	BUILTIN = "BUILTIN"
	// ARRAY object type
	ARRAY = "ARRAY"
	// HASH object type
	HASH = "HASH"
)

// Object interface is implemented by the objects.
type Object interface {
	Inspect() string
	Type() Type
}

// BuiltinFunction is a built-in function type.
type BuiltinFunction func(args ...Object) Object

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

// Void object.
type Void struct{}

// Inspect returns void.
func (v *Void) Inspect() string {
	return "null"
}

// Type returns the void object type.
func (v *Void) Type() Type {
	return VOID
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

// Builtin is a wrapper over built-in function.
type Builtin struct {
	Fn BuiltinFunction
}

// Type returns the built-ins' type
func (b *Builtin) Type() Type {
	return BUILTIN
}

// Inspect returns the builtin function representation
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// Array represents slice of objects
type Array struct {
	Elements []Object
}

// Type returns array type
func (a *Array) Type() Type {
	return ARRAY
}

// Inspect returns stringified array
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Hashable interface represents types that can be keys in hash object.
type Hashable interface {
	HashKey() HashKey
}

// HashKey is key in Hash.
type HashKey struct {
	Type  Type
	Value uint64
}

// HashKey returns HashKey created from a Boolean.
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// HashKey returns HashKey created from a Integer.
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// HashKey returns HashKey created from a String.
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// HashPair represents
type HashPair struct {
	Key   Object
	Value Object
}

// Hash represents the Hash Object Type.
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type returns the Hash object type.
func (h *Hash) Type() Type {
	return HASH
}

// Inspect returns stringified Hash object.
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
