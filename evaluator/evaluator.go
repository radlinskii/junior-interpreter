package evaluator

import (
	"github.com/radlinskii/interpreter/ast"
	"github.com/radlinskii/interpreter/object"
)

var (
	// TRUE is a single object that all the appeareances of boolean nodes with value "true" will point to.
	TRUE = &object.Boolean{Value: true}
	// FALSE is a single object that all the appeareances of boolean nodes with value "false" will point to.
	FALSE = &object.Boolean{Value: false}
	// NULL is a single object that all the appeareances of nodes without a value will point to.
	NULL = &object.Null{}
)

// Eval evaluates the program
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	//Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		if node.Value {
			return TRUE
		}
		return FALSE
	}

	return nil
}

func evalStatements(stmnts []ast.Statement) object.Object {
	var result object.Object

	for _, stmnt := range stmnts {
		result = Eval(stmnt)
	}

	return result
}
