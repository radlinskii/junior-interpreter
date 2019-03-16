package parser

import (
	"fmt"
	"testing"

	"github.com/radlinskii/interpreter/ast"
	"github.com/radlinskii/interpreter/lexer"
)

func testParsingInput(t *testing.T, input string, stmntLen int) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != stmntLen {
		t.Fatalf("program.Statements doesn't contain %d statements. got = %d", stmntLen, len(program.Statements))
	}

	return program
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser encountered %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"var x = 5;", "x", 5},
		{"var y = true;", "y", true},
		{"var z = y;", "z", "y"},
	}

	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)

		stmnt := program.Statements[0]
		if !testVarStatement(t, stmnt, tt.expectedIdentifier) {
			return
		}

		val := stmnt.(*ast.VarStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got='%q'", s.TokenLiteral())
		return false
	}

	varStmnt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("varStmnt not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmnt.Name.Value != name {
		t.Errorf("varStmnt.Name.Value not '%s'. got=%s", name, varStmnt.Name.Value)
		return false
	}

	if varStmnt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not %s. got=%s", name, varStmnt.Name)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return y;", "y"},
	}

	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)
		returnStmnt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmnt not *ast.ReturnStatement. got=%T", program.Statements[0])
			continue
		}
		if returnStmnt.TokenLiteral() != "return" {
			t.Errorf("returnStmnt.TokenLiteral() not `return`. got=%q", returnStmnt.TokenLiteral())
		}

		if returnStmnt.ReturnValue.TokenLiteral() != fmt.Sprintf("%v", tt.expectedValue) {
			t.Errorf("returnStmnt.ReturnValue expected=%v. got=%v", tt.expectedValue, returnStmnt.ReturnValue)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"foo;", "foo"},
		{"bar;", "bar"},
		{"baz", "baz"},
	}
	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)

		stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%q", program.Statements[0])
		}

		if !testIdentifier(t, stmnt.Expression, tt.expected) {
			return
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5;", 5},
		{"15", 15},
	}
	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)

		stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%q", program.Statements[0])
		}

		if !testIntegerLiteral(t, stmnt.Expression, tt.expected) {
			return
		}
	}
}

func TestBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false", false},
	}
	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)

		stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%q", program.Statements[0])
		}

		if !testBooleanLiteral(t, stmnt.Expression, tt.expected) {
			return
		}
	}
}
func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)

		stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%q", program.Statements[0])
		}

		if !testPrefixExpression(t, stmnt.Expression, tt.operator, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5+5;", 5, "+", 5},
		{"5-5;", 5, "-", 5},
		{"5*5;", 5, "*", 5},
		{"5/5;", 5, "/", 5},
		{"5>5;", 5, ">", 5},
		{"5<5;", 5, "<", 5},
		{"5>=5;", 5, ">=", 5},
		{"5<=5;", 5, "<=", 5},
		{"5==5;", 5, "==", 5},
		{"5!=5;", 5, "!=", 5},
		{"true==true;", true, "==", true},
		{"false==false;", false, "==", false},
		{"true!=false;", true, "!=", false},
	}

	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)

		stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%q", program.Statements[0])
		}

		if !testInfixExpression(t, stmnt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true;", "true"},
		{"false;", "false"},
		{"5+2 == 7;", "((5 + 2) == 7)"},
		{" 5+2 == 7 == true;", "(((5 + 2) == 7) == true)"},
		{"-a * b", "((-a) * b)"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a + -b", "(a + (-b))"},
		{"a * b + c", "((a * b) + c)"},
		{"a + b / c", "(a + (b / c))"},
		{"5 > 4 == 2 < 3", "((5 > 4) == (2 < 3))"},
		{"5 * 4 > 2 / 3 + 1", "((5 * 4) > ((2 / 3) + 1))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"10 / (5 + 5)", "(10 / (5 + 5))"},
		{"-(2 + 3)", "(-(2 + 3))"},
		{"!(true  == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(1,2+3,add(4,5))", "add(1, (2 + 3), add(4, 5))"},
		{"add(a+b+c*d/f, g)", "add(((a + b) + ((c * d) / f)), g)"},
	}

	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `
	if (x < y) {
		x;
	}`

	program := testParsingInput(t, input, 1)

	stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.Expression. got=%T", program.Statements[0])
	}

	exp, ok := stmnt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmnt.Expression is not *ast.IfExpression. got=%T", stmnt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative was not nil. got=%+v", exp.Alternative)
	}

}

func TestIfElseExpression(t *testing.T) {
	input := `
	if (x < y) {
		x;
	} else {
		y;
	}`

	program := testParsingInput(t, input, 1)

	stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.Expression. got=%T", program.Statements[0])
	}

	exp, ok := stmnt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmnt.Expression is not *ast.IfExpression. got=%T", stmnt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Alternative.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right interface{}) bool {
	pe, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("expressionStatement is not *ast.PrefixExpression. got=%q", exp)
		return false
	}

	if pe.Operator != operator {
		t.Fatalf("exp.Operator is not %s. got=%s", operator, pe.Operator)
		return false
	}

	if !testLiteralExpression(t, pe.Right, right) {
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	ie, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not *ast.InfixExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, ie.Left, left) {
		return false
	}

	if ie.Operator != operator {
		t.Errorf("exp.Operator is not %s. got=%q", operator, ie.Operator)
		return false
	}

	if !testLiteralExpression(t, ie.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integer, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp is not *ast.IntegerLiteral. got=%q", exp)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() not %d. got=%s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not *ast.Identifier. got=%q", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp is not *ast.BooleanLiteral. got=%q", exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() not %t. got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func TestFunctionLiteral(t *testing.T) {
	input := `fun(x,y) {x + y;}`

	program := testParsingInput(t, input, 1)

	stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmnt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("smnt.Expression is not *ast.FunctionalExpression. got=%T", stmnt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("expected 2 function parameters. got=%d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements length expected to be 1. got=%d", len(function.Body.Statements))
	}

	bodyStmnt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body statement is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmnt.Expression, "x", "+", "y")
}

func TestFunctionParametersParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fun() {};", expectedParams: []string{}},
		{input: "fun(x) {};", expectedParams: []string{"x"}},
		{input: "fun(x,y) {};", expectedParams: []string{"x", "y"}},
		{input: "fun(x,y,z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)

		stmnt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmnt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("parameters length wrong, want %d. got=%d", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2*3, 4+5);"

	program := testParsingInput(t, input, 1)

	stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ce, ok := stmnt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("smnt.Expression is not *ast.CallExpression. got=%T", stmnt.Expression)
	}

	if !testIdentifier(t, ce.Function, "add") {
		return
	}

	if len(ce.Arguments) != 3 {
		t.Fatalf("wrong length of arguments, expected 3. got=%d", len(ce.Arguments))
	}

	testLiteralExpression(t, ce.Arguments[0], 1)
	testInfixExpression(t, ce.Arguments[1], 2, "*", 3)
	testInfixExpression(t, ce.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	program := testParsingInput(t, input, 1)

	stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	literal, ok := stmnt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmnt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}
