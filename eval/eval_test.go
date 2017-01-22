package eval

import (
	"jacob/black/lexer"
	"jacob/black/object"
	"jacob/black/parser"
	"testing"
)

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to 'len' not supported, got 'int'"},
		{`len("one", "two")`, "wrong number of arguments. got '2', expected '1'"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.43", 5.43},
		{"10.0", 10.0},
		{"-10.3", -10.3},
		{"5.0 + 5.0 + 5.0 + 5.0 - 10.0", 10.0},
		{"2.0 * 2.0 * 2.0 * 2.0 * 2.0", 32.0},
		{"-50.0 + 100.0 + -50.0", 0.0},
		{"5.0 * 2 + 10.0", 20.0},
		{"5 + 2.0 * 10", 25.0},
		{"20 + 2 * -10.0", 0.0},
		{"50 / 2.0 * 2 + 10.0", 60.0},
		{"2 * (5 + 10.0)", 30.0},
		{"3 * 3.0 * 3 + 10", 37.0},
		{"3 * (3 * 3.0) + 10", 37.0},
		{"(5 + 10 * 2.0 + 15 / 3) * 2.0 + -10", 50.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1.0 < 2.0", true},
		{"1.1 > 2.4", false},
		{"1.0 < 1.0", false},
		{"1.0 > 1.0", false},
		{"1.13 == 1.13", true},
		{"1.13 != 1.13", false},
		{"1 == 2.0", false},
		{"1.1 != 2", true},
		{"0 == true", false},
		{"0 == false", false},
		{"1 == true", false},
		{"1 == false", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
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
		{"!0", true},
		{"!0.0", true},
		{"!0.1", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", 10},
		{"if false { 10 }", nil},
		{"if 1 { 10 }", 10},
		{"if 1 < 2 { 10 }", 10},
		{"if 1 > 2 { 10 }", nil},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 < 2 { 10 } else { 20 }", 10},
		{"if 0 { 10 } else { 5 }", 5},
		{"if 1 { 10 } else { 5 }", 10},
		{"if !0 { 10 } else { 5 }", 10},
		{"if !1 { 10 } else { 5 }", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"ret 10;", 10},
		{"ret 10; 9;", 10},
		{"ret 2 * 5; 9;", 10},
		{"9; ret 2 * 5; 9;", 10},
		{"if 10 > 1 { ret 10; }", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    ret 10;
  }

  ret 1;
}
`,
			10,
		},
		{
			`
		let f = |x| {
		  ret x;
		  x + 10;
		};
		f(10);`,
			10,
		},
		{
			`
		let f = |x| {
		   let result = x + 10;
		   ret result;
		   ret 10;
		};
		f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"cannot apply operator '+' for type 'int' and 'bool'",
		},
		{
			"5 + true; 5;",
			"cannot apply operator '+' for type 'int' and 'bool'",
		},
		{
			"-true",
			"unknown operator '-' for type 'bool'",
		},
		{
			"true + false;",
			"cannot apply operator '+' for type 'bool' and 'bool'",
		},
		{
			"true + false + true + false;",
			"cannot apply operator '+' for type 'bool' and 'bool'",
		},
		{
			"5; true + false; 5",
			"cannot apply operator '+' for type 'bool' and 'bool'",
		},
		{
			"if 10 > 1 { true + false; }",
			"cannot apply operator '+' for type 'bool' and 'bool'",
		},
		{
			`
		if 10 > 1 {
		  if 10 > 1 {
		    ret true + false;
		  }

		  ret 1;
		}
		`,
			"cannot apply operator '+' for type 'bool' and 'bool'",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"cannot apply operator '-' for type 'string' and 'string'",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "|x| x + 2;"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Params) != 1 {
		t.Fatalf("function has wrong parameters. Params=%+v",
			fn.Params)
	}

	if fn.Params[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Params[0])
	}

	expectedBody := "| (x + 2) "

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = |x| { x; }; identity(5);", 5},
		{"let identity = |x| { ret x; }; identity(5);", 5},
		{"let double = |x| { x * 2; }; double(5);", 10},
		{"let add = |x, y| { x + y; }; add(5, 5);", 10},
		{"let add = |x, y| { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"|x| { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
let first = 10;
let second = 10;
let third = 10;

let ourFunction = |first| {
  let second = 20;

  first + second + third;
};

ourFunction(20) + first + second;`

	testIntegerObject(t, testEval(input), 70)
}

func TestClassAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`let person = || {
			let age = 5
			ret || person
		};
		let p = person!
		p.age`, 5},
		{`let person = || {
			let age = 5
			ret || person
		};
		let house = || {
			let tennant = person!
			ret || house
		}
		let h = house!
		h.tennant.age`, 5},
		{`let person = || {
			let age = 5
			ret || person
		};
		let p = person!
		p.age = 6
		p.age`, 6},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.WithString(input, "testeval")
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f",
			result.Value, expected)
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
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testNilObject(t *testing.T, obj object.Object) bool {
	if obj != ConstNil {
		t.Errorf("object is not nil. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
