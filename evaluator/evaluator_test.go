package evaluator

import (
	"testing"

	"github.com/radlinskii/interpreter/lexer"
	"github.com/radlinskii/interpreter/object"
	"github.com/radlinskii/interpreter/parser"
)

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) != 0 {

		t.Errorf("parser encountered %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}

		t.FailNow()
	}

	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5;", 5},
		{"10;", 10},
		{"-55;", -55},
		{"-10;", -10},
		{"5 + 5 + 5 - 10;", 5},
		{"-10 + 23;", 13},
		{"2*2*2*2;", 16},
		{"2 + 3 * 4;", 14},
		{"2 * 3 + 4;", 10},
		{"2 + 3 * -4;", -10},
		{"50 / 2 * 2 + 10;", 60},
		{"2 * (5 + 10);", 30},
		{"3 * 3 * 3 + 10;", 37},
		{"3 * (3 * 3) + 10;", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10;", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
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
		{"1 < 2;", true},
		{"1 > 2;", false},
		{"1 < 1;", false},
		{"1 > 1;", false},
		{"1 >= 1;", true},
		{"1 <= 1;", true},
		{"1 == 1;", true},
		{"1 != 1;", false},
		{"1 != 2;", true},
		{"1 == 2;", false},
		{"true == true;", true},
		{"false == true;", false},
		{"false == false;", true},
		{"false != true;", true},
		{"(1 < 2) == true;", true},
		{"(1 > 2) != true;", true},
		{"(1 < 2) == false;", false},
		{"(1 < 2) == (4 >= 8);", false},
		{"(1 < 2) != (4 >= 8);", true},
		{"1 < 2 != 4 >= 8;", true},
		{"1 < 2 != 4 <= 8;", false},
		{"1 < 2 == 4 <= 8;", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
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
		{"!true;", false},
		{"!false;", true},
		{"!!true;", true},
		{"!!false;", false},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			return
		}
	}
}

func TestIfElseStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10; }", 10},
		{"if (false) { 10; }", nil},
		{"if (true) { 10; }", 10},
		{"if (1 < 2) { 10; }", 10},
		{"if (1 > 2) { 10; }", nil},
		{"if (1 > 2) { 10; } else { 20; }", 20},
		{"if (1 < 2) { 10; } else { 20; }", 10},
		{`
			if (true) {
				if (true) {
					100;
				}
				20;
			} else {
				30;
			}
		`, 20},
		{`
			if (true) {
				if (false) {
					10;
				}

				20;
			} else {
				30;
			}
		`, 20},
		{`
			if (false) {
				if (false) {
					10;
				}

				20;
			} else {
				30;
			}
		`, 30},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
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
		{
			`
			const a = fun(x) {
				if (x <= 10) {
					return x;
				}
				return 5;
			};

			a(10);`,
			10,
		},
		{
			`
			const a = fun(x) {
				return x*5;
			};

			a(5);`,
			25,
		},
		{
			`
			const a = fun(x) {
				return 10;
				return 5;
			};

			a(10);`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
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
		{`"Hell" - "world";`, "unknown operator: STRING - STRING"},
		{`5 + "worlds";`, "type mismatch: INTEGER + STRING"},
		{`{fun(x) { return x +1; }: "Monkey"}[fun(x) { return x +1; }];`, "FUNCTION can't be used as hash key"},
		{`{"key": "Monkey"}[fun(x) { return x +1; }];`, "index operator not supported: HASH[FUNCTION]"},
		{`
			const a = 5;
			return a;`,
			"return statement not perrmitted outside function body",
		},
		{`
			const a = 5;
			if (a < 10) {
				return a;
			}`,
			"return statement not perrmitted outside function body",
		},
		{`
			const a = fun(x) {
				if (x < 10) {
					return "dupa";
				}
				return "dupa";
			};

			return a(3);`,
			"return statement not perrmitted outside function body",
		},
		{`
			const a = fun(x) {
				if (x < 10) {
					return "duupa";
				}
				15;
			};

			a(20);`,
			"missing return at the end of function body",
		},
		{`
			if (1 < 2) {
				const foobar = "baaaz";
			}

			print(foobar);`,
			"unknown identifier: foobar"},
		{`
			const foobar = "foo";
			if (1 < 2) {
				const foobar = "bar";
			}

			const foobar = "baz";`,
			`redeclared constant: "foobar" in one block`},
		{`
			const foobar = "foo";
			if (1 < 2) {
				const fizz = "bar";
				const fizz = "baz";
			}

			print(foobar);`,
			`redeclared constant: "fizz" in one block`},
		{`
			const someFunc = fun(x) {
				const x = "oh no!";

				return;
			};

			print(someFunc("oh yes"));`,
			`redeclared constant: "x" in one block`},
		{`
			if (1) {
				print("1 is truthy??");
			} else {
				print("there is no 'truthy'! ");
			}`,
			`expected BOOLEAN as condition in if-statement got: INTEGER`},

		{`
			if (!1) {
				print("!1 is falsy??");
			} else {
				print("there is no 'falsy'! ");
			}`,
			`expected BOOLEAN in negation expression, got: INTEGER`},
	}

	for _, tt := range tests {
		if !testErrorObject(t, testEval(t, tt.input), tt.expectedMsg) {
			return
		}
	}
}

func testErrorObject(t *testing.T, obj object.Object, expectedMessage string) bool {
	errObj, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object is not an Error. got=%T(%+v)", obj, obj)
		return false
	}

	if errObj.Message != expectedMessage {
		t.Errorf("wrong Error Message. expected=%q, got %q", expectedMessage, errObj.Message)
		return false
	}

	return true
}

func TestConstStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"const a = 5; a;", 5},
		{"const a = 5 * 5; a;", 25},
		{"const a = 5; const b = a; b;", 5},
		{"const a = 5; const b = a; const c = a + b + 5; c + 5;", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		if !testIntegerObject(t, evaluated, tt.expected) {
			return
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fun(x) { return x + 2; };"
	expectedBody := "return (x + 2);"

	evaluated := testEval(t, input)
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
		{"const identity = fun(x) { return x; }; identity(5);", 5},
		{"const identity = fun(x) { return x; }; identity(5);", 5},
		{"const double = fun(x) { return x * 2; }; double(5);", 10},
		{"const add = fun(x, y) { return x + y; }; add(5, 10);", 15},
		{"const add = fun(x, y) { return x + y; }; add(5, add(6, 4));", 15},
		{"fun(x) { return x; }(99);", 99},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEval(t, tt.input), tt.expected) {
			return
		}
	}
}

func TestClosures(t *testing.T) {
	input := `
		const newAdder = fun(x) {
			return fun(y) {
				return x + y;
			};
		};

		const addFive = newAdder(5);
		addFive(5);
	`
	expected := 10

	if !testIntegerObject(t, testEval(t, input), int64(expected)) {
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

	if !testStringObject(t, testEval(t, input), "Hello World") {
		return
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " World";`

	if !testStringObject(t, testEval(t, input), "Hello World") {
		return
	}
}

func TestStringComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"hello" == "hello";`, true},
		{`"hello" != "hello";`, false},
		{`"hello" == " hello";`, false},
		{`"hello" != " hello";`, true},
		{`"hello" == "worlds";`, false},
	}

	for _, tt := range tests {
		if !testBooleanObject(t, testEval(t, tt.input), tt.expected) {
			return
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("");`, 0},
		{`len("four");`, 4},
		{`len("hello world");`, 11},
		{`len(1);`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two");`, "wrong number of arguments. got=2 want=1"},
		{`len([1,2,3,4]);`, 4},
		{`len([]);`, 0},
		{`first([1,2,3,4]);`, 1},
		{`first([]);`, nil},
		{`last([1,2,3,4]);`, 4},
		{`last([]);`, nil},
		{`rest([1,2,3,4]);`, []int{2, 3, 4}},
		{`rest([1]);`, []int{}},
		{`rest([]);`, nil},
		{`push([], 1);`, []int{1}},
		{`push([1,2],3);`, []int{1, 2, 3}},
		{`push([1,2,3]);`, "wrong number of arguments. got=1 want=2"},
		{`push([1,2,3],3,3);`, "wrong number of arguments. got=3 want=2"},
		{`push(true,3);`, "first argument to `push` not supported, got BOOLEAN"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch expected := tt.expected.(type) {
		case int:
			if !testIntegerObject(t, evaluated, int64(expected)) {
				return
			}
		case string:
			if !testErrorObject(t, evaluated, expected) {
				return
			}
		case []int:
			expectedArr := tt.expected.([]int)
			evaluatedArr := evaluated.(*object.Array).Elements
			if len(evaluated.(*object.Array).Elements) != len(expectedArr) {
				t.Errorf("evaluatedArr length expected to be %d, got %d", len(expectedArr), len(evaluatedArr))
			}
			for i, o := range evaluatedArr {
				if !testIntegerObject(t, o, int64(expectedArr[i])) {
					return
				}
			}
		case nil:
			if !testNullObject(t, evaluated) {
				return
			}
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, true, "word"];`

	evaluated := testEval(t, input)

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
		{"[1,2,3][0];", 1},
		{"[1,2,3][1];", 2},
		{"[1,2,3][2];", 3},
		{"const i = 0; [1][i];", 1},
		{"[1,2,3][1 + 1];", 3},
		{"const myArray = [1, 2, 3]; myArray[0];", 1},
		{"const myArray = [1, 2, 3]; myArray[0] + myArray[2];", 4},
		{"[1, 2, 3][-1];", "index out of boundaries"},
		{"[1, 2, 3][3];", "index out of boundaries"},
		{"[1, 2, 3][true];", "index operator not supported: ARRAY[BOOLEAN]"},
		{"54[1];", "index operator not supported: INTEGER[INTEGER]"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testErrorObject(t, evaluated, tt.expected.(string))
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `
	const two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6,
	};
	`
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 1}["foo"];`, 1},
		{`{"foo": 5}["bar"];`, "No hash pair in \"{foo: 5}\" with key \"bar\""},
		{`const key = "foo"; {"foo": 5}[key];`, 5},
		{`{}["foo"];`, "No hash pair in \"{}\" with key \"foo\""},
		{`{5: 10}[5];`, 10},
		{`{true: 5}[true];`, 5},
		{`{false: 5}[false];`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testErrorObject(t, evaluated, tt.expected.(string))
		}
	}
}

func TestVoidFunction(t *testing.T) {
	input := `
	const foo = fun(x) {
		const b  = x + 2;

		print(b);

		return;
	};

	foo(4);
	`

	evaluated := testEval(t, input)
	if evaluated != VOID {
		t.Errorf("object is not NULL. got=%T", evaluated)
	}
}

func TestBlockScope(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
		const a = 10;
		if (1 < 10) {
			const a = 5;
		} else {
			const a = 15;
		}

		a;
		`, 10},
		{`
		const a = 10;
		const foo = fun(x) {
			const a = x + x;
			return;
		};

		foo(a);
		a;
		`, 10},
		{`
		const a = 10;
		const foo = fun(x) {
			const a = x / 2;

			return a;
		};

		foo(a);
		`, 5},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEval(t, tt.input), tt.expected) {
			return
		}
	}

}
