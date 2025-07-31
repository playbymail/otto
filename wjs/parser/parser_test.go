// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package parser

import (
	"testing"

	"github.com/playbymail/otto/wjs/ast"
	"github.com/playbymail/otto/wjs/lexer"
)

func TestLetStatement(t *testing.T) {
	input := "let x = 5;"
	
	l := lexer.New("test", input)
	tokens := l.AllTokens()
	
	// Debug: print tokens
	t.Logf("Tokens:")
	for i, token := range tokens {
		t.Logf("  %d: %s", i, token)
	}
	
	p := New(tokens)
	program := p.ParseProgram()
	
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	
	if len(program.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Stmts))
	}
	
	t.Logf("Statement type: %T", program.Stmts[0])
	
	stmt, ok := program.Stmts[0].(*ast.LetStmt)
	if !ok {
		t.Fatalf("Expected LetStmt, got %T", program.Stmts[0])
	}
	
	if stmt.Name.Name != "x" {
		t.Errorf("Expected name 'x', got %s", stmt.Name.Name)
	}
}

func TestBinaryExpression(t *testing.T) {
	input := "5 + 3;"
	
	l := lexer.New("test", input)
	tokens := l.AllTokens()
	
	p := New(tokens)
	program := p.ParseProgram()
	
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	
	if len(program.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Stmts))
	}
	
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Expected ExprStmt, got %T", program.Stmts[0])
	}
	
	expr, ok := stmt.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Expected BinaryExpr, got %T", stmt.Value)
	}
	
	if expr.Operator != "+" {
		t.Errorf("Expected operator '+', got %s", expr.Operator)
	}
}

func TestCallExpression(t *testing.T) {
	input := "print(42);"
	
	l := lexer.New("test", input)
	tokens := l.AllTokens()
	
	p := New(tokens)
	program := p.ParseProgram()
	
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	
	if len(program.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Stmts))
	}
	
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Expected ExprStmt, got %T", program.Stmts[0])
	}
	
	call, ok := stmt.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("Expected CallExpr, got %T", stmt.Value)
	}
	
	ident, ok := call.Callee.(*ast.Ident)
	if !ok {
		t.Fatalf("Expected Ident callee, got %T", call.Callee)
	}
	
	if ident.Name != "print" {
		t.Errorf("Expected function name 'print', got %s", ident.Name)
	}
	
	if len(call.Args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(call.Args))
	}
}

func TestAssignmentStatement(t *testing.T) {
	input := "x = 10;"
	
	l := lexer.New("test", input)
	tokens := l.AllTokens()
	
	p := New(tokens)
	program := p.ParseProgram()
	
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	
	if len(program.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Stmts))
	}
	
	stmt, ok := program.Stmts[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("Expected AssignStmt, got %T", program.Stmts[0])
	}
	
	ident, ok := stmt.Target.(*ast.Ident)
	if !ok {
		t.Fatalf("Expected Ident target, got %T", stmt.Target)
	}
	
	if ident.Name != "x" {
		t.Errorf("Expected target 'x', got %s", ident.Name)
	}
}
