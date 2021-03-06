package evaluator

import (
	"bytes"
	"fmt"

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
	// VOID is a single object that all the appeareances of nodes without a value will point to.
	VOID = &object.Void{}
)

var programOutput bytes.Buffer

// eval evaluates the AST
func eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return eval(node.Expression, env)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.ConstStatement:
		return evalConstStatement(node, env)
	//Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return evalBoolToBooleanObjectReference(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.PrefixExpression:
		right := eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		fun := eval(node.Function, env)
		if isError(fun) {
			return fun
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fun, args)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalIndexExpression(left, right)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	default:
		return nil
	}
}

// EvalProgram starts evaluation of the AST.
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmnt := range program.Statements {
		result = eval(stmnt, env)

		switch result := result.(type) {
		case *object.Return:
			return newError("return statement not permitted outside function body")
		case *object.Error:
			return result
		}
	}

	return result
}

// EvalProgram starts evaluation of the AST.
func EvalProgram(program *ast.Program, env *object.Environment) string {
	evaluated := evalProgram(program, env)

	programOutput.WriteString(evaluated.Inspect())

	retStr := programOutput.String()
	programOutput.Reset()

	return retStr
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	blockEnv := object.NewEnclosedEnvironment(env)

	for _, stmnt := range block.Statements {
		result = eval(stmnt, blockEnv)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN || rt == object.ERROR {
				return result
			}
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
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return newError("expected BOOLEAN in negation expression, got: %s", right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type(): // handling type mismatch error first
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return evalBoolToBooleanObjectReference(left == right)
	case operator == "!=":
		return evalBoolToBooleanObjectReference(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return evalBoolToBooleanObjectReference(leftVal == rightVal)
	case "!=":
		return evalBoolToBooleanObjectReference(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBoolToBooleanObjectReference(val bool) object.Object {
	if val {
		return TRUE
	}
	return FALSE
}

func evalIfStatement(ie *ast.IfStatement, env *object.Environment) object.Object {
	condition := eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	isConditionTrue, ok := isTruthy(condition)
	if !ok {
		return newError("expected BOOLEAN as condition in if-statement got: %s", condition.Type())
	}

	if isConditionTrue {
		return eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return eval(ie.Alternative, env)
	}

	return NULL
}

func isTruthy(obj object.Object) (val, ok bool) {
	switch obj {
	case FALSE:
		val = false
		ok = true
	case TRUE:
		val = true
		ok = true
	default:
		val = false
		ok = false
	}

	return
}

func evalIdentifier(i *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(i.Value); ok {
		return val
	}

	if builtin, ok := builtins[i.Value]; ok {
		return builtin
	}

	return newError("unknown identifier: %s", i.Value)
}

func evalIndexExpression(left, right object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && right.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, right)
	case left.Type() == object.HASH:
		return evalHashIndexExpression(left, right)
	default:
		return newError("index operator not supported: %s[%s]", left.Type(), right.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if i < 0 || i > max {
		return newError("index out of boundaries")
	}

	return arrayObject.Elements[i]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("index operator not supported: %s[%s]", hash.Type(), index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return newError("No hash pair in %q with key %q", hash.Inspect(), index.Inspect())
	}

	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("%s can't be used as hash key", key.Type())
		}

		value := eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalReturnStatement(rs *ast.ReturnStatement, env *object.Environment) object.Object {
	if rs.ReturnValue == nil {
		return VOID
	}
	val := eval(rs.ReturnValue, env)
	if isError(val) {
		return val
	}
	return &object.Return{Value: val}
}

func evalConstStatement(cs *ast.ConstStatement, env *object.Environment) object.Object {
	if _, ok := env.ShallowGet(cs.Name.Value); ok {
		return newError("redeclared constant: %q in one block", cs.Name.Value)
	}

	val := eval(cs.Value, env)
	if isError(val) {
		return val
	}

	return env.Set(cs.Name.Value, val)
}

func applyFunction(fun object.Object, args []object.Object) object.Object {
	switch function := fun.(type) {
	case *object.Function:
		extendedEnv := extendedFunctionEnv(function, args)
		evaluated := evalFunctionBody(function.Body, extendedEnv)

		if isError(evaluated) {
			return evaluated
		}

		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return function.Fn(args...)
	default:
		return newError("not a function: %s", function.Type())
	}
}

func evalFunctionBody(body *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmnt := range body.Statements {
		result = eval(stmnt, env)

		if result != nil {
			rt := result.Type()
			if result == VOID || rt == object.RETURN || rt == object.ERROR {
				return result
			}
		}
	}

	return newError("missing return at the end of function body")
}

func extendedFunctionEnv(fun *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fun.Env)

	for paramIdx, param := range fun.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if rtrn, ok := obj.(*object.Return); ok {
		return rtrn.Value
	}

	return obj
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}
