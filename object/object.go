package object

import (
	"fmt"
)

// Type represents different object types.
type Type string

const (
	// INTEGER object type
	INTEGER = "INTEGER"
	// BOOLEAN object type
	BOOLEAN = "BOOLEAN"
	// NULL object type
	NULL = "NULL"
	// RETURN object type
	RETURN = "RETURN"
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

// Return object.
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
