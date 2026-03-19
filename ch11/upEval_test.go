package main

import (
	"testing"
)

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
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
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

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
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

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
			`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
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

// -----------------------------------------------------------------------
// NEW TESTS FOR MUTABILITY AND LOOPS
// -----------------------------------------------------------------------

func TestAssignments(t *testing.T) {
	// This verifies that Environment.Update is working correctly
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a = 10; a;", 10},
		{"let a = 5; a = a + 10; a;", 15},
		{"let a = 5; let b = 10; a = b; a;", 10},
		{"let a = 5; a = a * 2; a = a + 5; a;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestForLoops(t *testing.T) {
	// This verifies the parsing AND execution of the new C-style loop
	tests := []struct {
		input    string
		expected int64
	}{
		// Standard counting loop
		{
			`let sum = 0;
			 for (let i = 0; i < 5; i = i + 1) {
				 sum = sum + 1;
			 }
			 sum;`,
			5,
		},
		// Using the loop index in calculation
		{
			`let sum = 0;
			 for (let i = 0; i < 4; i = i + 1) {
				 sum = sum + i;
			 }
			 sum;`,
			6, // 0 + 1 + 2 + 3
		},
		// Double nested loop
		{
			`let sum = 0;
			 for (let i = 0; i < 2; i = i + 1) {
				 for (let j = 0; j < 2; j = j + 1) {
					 sum = sum + 1;
				 }
			 }
			 sum;`,
			4,
		},
		// Loop with return (Early Exit)
		{
			`func() {
				for (let i = 0; i < 10; i = i + 1) {
					if (i == 5) {
						return i;
					}
				}
				return -1;
			 }();`,
			5,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosuresAndScoping(t *testing.T) {
	//
	// Verifies that Update() correctly walks up the scope chain
	input := `
	let newAdder = func(x) {
		func(y) { x + y };
	};
	let addTwo = newAdder(2);
	addTwo(2);
	`
	testIntegerObject(t, testEval(input), 4)
}

func TestOuterScopeMutation(t *testing.T) {
	// Verifies that a closure can modify a variable in the outer scope
	input := `
	let global = 10;
	let modifier = func() {
		global = 20;
	};
	modifier();
	global;
	`
	testIntegerObject(t, testEval(input), 20)
}

// -----------------------------------------------------------------------
// HELPERS
// -----------------------------------------------------------------------

func testEval(input string) Object {
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	env := NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj Object, expected int64) bool {
	result, ok := obj.(*Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	result, ok := obj.(*Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj Object) bool {
	if obj.Type() != NULL_OBJ {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
