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
	input := `
	var x = 5;
	var y = 10;
	var foo = 9999;
	`

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	program := testParsingInput(t, input, 3)

	for i, tt := range tests {
		stmnt := program.Statements[i]
		if !testVarStatement(t, stmnt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	t.Log(name)
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
	input := `
	return 5;
	return 10;
	return 9999;
	`

	program := testParsingInput(t, input, 3)

	for _, stmnt := range program.Statements {
		returnStmnt, ok := stmnt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmnt not *ast.ReturnStatement. got=%T", stmnt)
			continue
		}
		if returnStmnt.TokenLiteral() != "return" {
			t.Errorf("returnStmnt.TokenLiteral() not `return`. got=%q", returnStmnt.TokenLiteral())
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
	}

	for _, tt := range tests {
		program := testParsingInput(t, tt.input, 1)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
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
