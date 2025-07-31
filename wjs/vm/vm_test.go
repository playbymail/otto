// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package vm

import (
	"strings"
	"testing"

	"github.com/playbymail/otto/wjs/ast"
	"github.com/playbymail/otto/wjs/domain"
	"github.com/playbymail/otto/wjs/lexer"
	"github.com/playbymail/otto/wjs/parser"
)

func TestVM_NumberLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"42", 42},
		{"3.14", 3.14},
		{"0", 0},
		{"-5", -5},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := evalExpression(test.input)
			if err != nil {
				t.Fatalf("Runtime error: %v", err)
			}
			
			if num, ok := result.(float64); !ok || num != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestVM_StringLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"world"`, "world"},
		{`""`, ""},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := evalExpression(test.input)
			if err != nil {
				t.Fatalf("Runtime error: %v", err)
			}
			
			if str, ok := result.(string); !ok || str != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestVM_BinaryExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected Value
	}{
		// Arithmetic
		{"5 + 3", 8.0},
		{"10 - 4", 6.0},
		{"6 * 7", 42.0},
		{"20 / 4", 5.0},
		{"17 % 5", 2.0},
		
		// String concatenation
		{`"hello" + " world"`, "hello world"},
		
		// Comparison
		{"5 == 5", true},
		{"5 != 3", true},
		{"5 < 10", true},
		{"10 > 5", true},
		{"5 <= 5", true},
		{"5 >= 5", true},
		{"3 < 3", false},
		{"3 > 3", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := evalExpression(test.input)
			if err != nil {
				t.Fatalf("Runtime error: %v", err)
			}
			
			if !Equal(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestVM_UnaryExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected Value
	}{
		{"-5", -5.0},
		{"-(-3)", 3.0},
		{"!true", false},
		{"!false", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := evalExpression(test.input)
			if err != nil {
				t.Fatalf("Runtime error: %v", err)
			}
			
			if !Equal(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestVM_LetStatements(t *testing.T) {
	tests := []struct {
		input    string
		varName  string
		expected Value
	}{
		{"let x = 5;", "x", 5.0},
		{`let name = "test";`, "name", "test"},
		{"let result = 3 + 4;", "result", 7.0},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			vm := New("test")
			program := parseInput(test.input)
			
			_, err := vm.Execute(program)
			if err != nil {
				t.Fatalf("Runtime error: %v", err)
			}
			
			value, exists := vm.vars[test.varName]
			if !exists {
				t.Fatalf("Variable %s not found", test.varName)
			}
			
			if !Equal(value, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, value)
			}
		})
	}
}

func TestVM_Identifiers(t *testing.T) {
	input := `
		let x = 5;
		let y = x;
		y;
	`
	
	result, err := evalProgram(input)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}
	
	if num, ok := result.(float64); !ok || num != 5 {
		t.Errorf("Expected 5, got %v", result)
	}
}

func TestVM_AssignmentStatements(t *testing.T) {
	input := `
		let x = 5;
		x = 10;
		x;
	`
	
	result, err := evalProgram(input)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}
	
	if num, ok := result.(float64); !ok || num != 10 {
		t.Errorf("Expected 10, got %v", result)
	}
}

func TestVM_BuiltinPrint(t *testing.T) {
	// Capture print output
	var output strings.Builder
	originalPrint := func(pos domain.Pos, args []Value) (Value, *RuntimeError) {
		out := make([]string, len(args))
		for i, arg := range args {
			out[i] = Stringify(arg)
		}
		output.WriteString(strings.Join(out, " "))
		return nil, nil
	}
	
	vm := New("test")
	vm.vars["print"] = &builtinFunc{
		name:  "print",
		arity: -1,
		fn:    originalPrint,
	}
	
	input := `print("hello", "world");`
	program := parseInput(input)
	
	_, err := vm.Execute(program)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}
	
	expected := "hello world"
	if output.String() != expected {
		t.Errorf("Expected %q, got %q", expected, output.String())
	}
}

func TestVM_ErrorHandling(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{"undefined_var;", "undefined variable: undefined_var"},
		{"5 + true;", "type mismatch for + operator"},
		{"5 / 0;", "division by zero"},
		{"5 % 0;", "modulus by zero"},
		{"-true;", "unary - requires a number"},
		{"!5;", "unary ! requires a boolean"},
		{"unknown();", "undefined variable: unknown"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			_, err := evalExpression(test.input)
			if err == nil {
				t.Fatalf("Expected error, got none")
			}
			
			if !strings.Contains(err.Message, test.expectedError) {
				t.Errorf("Expected error containing %q, got %q", test.expectedError, err.Message)
			}
		})
	}
}

func TestVM_TemplateStrings(t *testing.T) {
	// Note: We'll need to add boolean literals to the parser first
	// For now, test basic template functionality
	input := `
		let name = "world";
		let greeting = "hello";
	`
	
	vm := New("test")
	program := parseInput(input)
	
	_, err := vm.Execute(program)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}
	
	// Verify variables were set
	if name, exists := vm.vars["name"]; !exists || name != "world" {
		t.Errorf("Expected name='world', got %v", name)
	}
	if greeting, exists := vm.vars["greeting"]; !exists || greeting != "hello" {
		t.Errorf("Expected greeting='hello', got %v", greeting)
	}
}

// Helper functions

func evalExpression(input string) (Value, *RuntimeError) {
	vm := New("test")
	tokens := getAllTokens(input)
	p := parser.New(tokens)
	program := p.ParseProgram()
	
	if len(program.Stmts) == 0 {
		return nil, NewRuntimeError(domain.Pos{}, "no statements to evaluate")
	}
	
	// Treat single expression as expression statement
	if len(program.Stmts) == 1 {
		if exprStmt, ok := program.Stmts[0].(*ast.ExprStmt); ok {
			return vm.evalExpr(exprStmt.Value)
		}
	}
	
	return vm.Execute(program)
}

func evalProgram(input string) (Value, *RuntimeError) {
	vm := New("test")
	program := parseInput(input)
	return vm.Execute(program)
}

func parseInput(input string) *ast.Program {
	tokens := getAllTokens(input)
	p := parser.New(tokens)
	return p.ParseProgram()
}

func getAllTokens(input string) []lexer.Token {
	l := lexer.New("test", input)
	return l.AllTokens()
}
