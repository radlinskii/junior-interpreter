package evaluator

import (
	"testing"

	"github.com/radlinskii/interpreter/lexer"
	"github.com/radlinskii/interpreter/object"
	"github.com/radlinskii/interpreter/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-55", -55},
		{"-10", -10},
		{"5 + 5 + 5 - 10", 5},
		{"-10 + 23", 13},
		{"2*2*2*2", 16},
		{"2 + 3 * 4", 14},
		{"2 * 3 + 4", 10},
		{"2 + 3 * -4", -10},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testIntegerObject(t, evaluated, tt.expected) {
			return
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 >= 1", true},
		{"1 <= 1", true},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 != 2", true},
		{"1 == 2", false},
		{"true == true", true},
		{"false == true", false},
		{"false == false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) != true", true},
		{"(1 < 2) == false", false},
		{"(1 < 2) == (4 >= 8)", false},
		{"(1 < 2) != (4 >= 8)", true},
		{"1 < 2 != 4 >= 8", true},
		{"1 < 2 != 4 <= 8", false},
		{"1 < 2 == 4 <= 8", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			return
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			return
		}
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10; }", 10},
		{"if (false) { 10; }", nil},
		{"if (1) { 10; }", 10},
		{"if (1 < 2) { 10; }", 10},
		{"if (1 > 2) { 10; }", nil},
		{"if (1 > 2) { 10; } else { 20; }", 20},
		{"if (1 < 2) { 10; } else { 20; }", 10},
		{`
			if (true) {
				if (true) {
					return 10;
				}

				return 20;
			} else {
				return 30;
			}
		`, 10},
		{`
			if (true) {
				if (false) {
					return 10;
				}

				return 20;
			} else {
				return 30;
			}
		`, 20},
		{`
			if (false) {
				if (false) {
					return 10;
				}

				return 20;
			} else {
				return 30;
			}
		`, 30},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if !ok {
			if !testNullObject(t, evaluated) {
				return
			}
		} else {
			if !testIntegerObject(t, evaluated, int64(integer)) {
				return
			}
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"2; return 10;", 10},
		{"true; return 10;", 10},
		{"return 10; 99;", 10},
		{"return 2 * 5 + 2; 99;", 12},
		{"true; return 2 + 5; false;", 7},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		expectedMsg string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 10;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if(10 > 1) { return true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar;", "unknown identifier: foobar"},
		{`"Hell" - "world"`, "unknown operator: STRING - STRING"},
		{`5 + "worlds";`, "type mismatch: INTEGER + STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("object is not Error. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMsg {
			t.Errorf("wrong Error Message. expected=%q, got %q", tt.expectedMsg, errObj.Message)
		}
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c + 5;", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if !testIntegerObject(t, evaluated, tt.expected) {
			return
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fun(x) { return x + 2; }"
	expectedBody := "return (x + 2);"

	evaluated := testEval(input)
	fun, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not a Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fun.Parameters) != 1 {
		t.Fatalf("function has wrong number of parameters. Parameters=%+v", fun.Parameters)
	}

	if fun.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fun.Parameters[0])
	}

	if fun.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fun.Body.String())
	}
}

func TestFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var identity = fun(x) { x; } identity(5);", 5},
		{"var identity = fun(x) { return x; } identity(5);", 5},
		{"var double = fun(x) { return x * 2; } double(5);", 10},
		{"var add = fun(x, y) { return x + y; } add(5, 10);", 15},
		{"var add = fun(x, y) { return x + y; } add(5, add(6, 4));", 15},
		{"fun(x) { return x; }(99);", 99},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEval(tt.input), tt.expected) {
			return
		}
	}
}

func TestClosures(t *testing.T) {
	input := `
		var newAdder = fun(x) {
			return fun(y) {
				return x + y;
			};
		};

		var addFive = newAdder(5);
		addFive(5);
	`
	expected := 10

	if !testIntegerObject(t, testEval(input), int64(expected)) {
		return
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Wrong Integer value, expected=%d, got=%d", expected, result.Value)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Wrong Boolean value, expected=%t, got=%t", expected, result.Value)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T", obj)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Wrong String value, expected=%s, got=%s", expected, result.Value)
		return false
	}
	return true
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World";`

	if !testStringObject(t, testEval(input), "Hello World") {
		return
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " World";`

	if !testStringObject(t, testEval(input), "Hello World") {
		return
	}
}

func TestStringComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"hello" == "hello"`, true},
		{`"hello" != "hello"`, false},
		{`"hello" == " hello"`, false},
		{`"hello" != " hello"`, true},
		{`"hello" == "worlds"`, false},
	}

	for _, tt := range tests {
		if !testBooleanObject(t, testEval(tt.input), tt.expected) {
			return
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2 want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			if !testIntegerObject(t, evaluated, int64(expected)) {
				return
			}
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, true, "word"];`

	evaluated := testEval(input)

	array, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(array.Elements) != 4 {
		t.Fatalf("array has wrong number of elements, expected=4, got=%d", len(array.Elements))
	}

	testIntegerObject(t, array.Elements[0], 1)
	testIntegerObject(t, array.Elements[1], 4)
	testBooleanObject(t, array.Elements[2], true)
	testStringObject(t, array.Elements[3], "word")
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1,2,3][0]", 1},
		{"[1,2,3][1]", 2},
		{"[1,2,3][2]", 3},
		{"var i = 0; [1][i];", 1},
		{"[1,2,3][1 + 1]", 3},
		{"var myArray = [1, 2, 3]; myArray[0]", 1},
		{"var myArray = [1, 2, 3]; myArray[0] + myArray[2]", 4},
		{"[1, 2, 3][-1]", nil}, // TODO this should definitely be an error
		{"[1, 2, 3][3]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
