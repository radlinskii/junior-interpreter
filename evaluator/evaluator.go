package evaluator

import (
	"github.com/radlinskii/interpreter/ast"
	"github.com/radlinskii/interpreter/object"
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
