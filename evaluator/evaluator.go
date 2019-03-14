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
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.Return{Value: val}
	//Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return evalBoolToBooleanObjectReference(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	default:
		return nil
	}
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmnt := range program.Statements {
		result = Eval(stmnt)

		if rtrn, ok := result.(*object.Return); ok {
			return rtrn.Value
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmnt := range block.Statements {
		result = Eval(stmnt)

		if result != nil && result.Type() == object.RETURN {
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL // TODO throw error??
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return NULL // TODO would be better to throw an error
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return evalBoolToBooleanObjectReference(left == right)
	case operator == "!=":
		return evalBoolToBooleanObjectReference(left != right)
	default:
		return NULL // TODO throw error??
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return evalBoolToBooleanObjectReference(leftVal < rightVal)
	case ">":
		return evalBoolToBooleanObjectReference(leftVal > rightVal)
	case "==":
		return evalBoolToBooleanObjectReference(leftVal == rightVal)
	case "!=":
		return evalBoolToBooleanObjectReference(leftVal != rightVal)
	case "<=":
		return evalBoolToBooleanObjectReference(leftVal <= rightVal)
	case ">=":
		return evalBoolToBooleanObjectReference(leftVal >= rightVal)
	default:
		return NULL // TODO definitely an error
	}
}

func evalBoolToBooleanObjectReference(val bool) object.Object {
	if val {
		return TRUE
	}
	return FALSE
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true // TODO ERROR non null type is used as condition
	}
}
