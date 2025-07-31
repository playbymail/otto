// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package vm provides the runtime for the WJS scripting language.
package vm

import (
	"fmt"
	"github.com/playbymail/otto/wjs/ast"
	"github.com/playbymail/otto/wjs/domain"
)

func New(script string) *VM {
	vm := &VM{
		vars:   map[string]Value{},
		script: script,
	}
	// Register built-in functions
	builtins := RegisterBuiltins(vm.defaultLoad, vm.defaultSave)
	for name, fn := range builtins {
		vm.vars[name] = fn
	}
	return vm
}

type VM struct {
	vars   map[string]Value // environment: variables and functions
	script string           // current script filename
}

// Execute runs the program and returns the last expression result (if any) and any runtime error.
func (vm *VM) Execute(program *ast.Program) (Value, *RuntimeError) {
	var lastValue Value
	
	for _, stmt := range program.Stmts {
		result, err := vm.evalStmt(stmt)
		if err != nil {
			return nil, err
		}
		if result != nil {
			lastValue = result
		}
	}
	
	return lastValue, nil
}

// evalStmt evaluates a statement and returns its value (if any) and runtime error (if any).
func (vm *VM) evalStmt(stmt ast.Stmt) (Value, *RuntimeError) {
	switch s := stmt.(type) {
	case *ast.LetStmt:
		return vm.evalLetStmt(s)
	case *ast.AssignStmt:
		return vm.evalAssignStmt(s)
	case *ast.ExprStmt:
		return vm.evalExprStmt(s)
	default:
		return nil, NewRuntimeError(s.Pos(), "unknown statement type: %T", s)
	}
}

// evalExpr evaluates an expression and returns its value and runtime error (if any).
func (vm *VM) evalExpr(expr ast.Expr) (Value, *RuntimeError) {
	switch e := expr.(type) {
	case *ast.NumberLit:
		return vm.evalNumberLit(e)
	case *ast.StringLit:
		return vm.evalStringLit(e)
	case *ast.Ident:
		return vm.evalIdent(e)
	case *ast.BinaryExpr:
		return vm.evalBinaryExpr(e)
	case *ast.UnaryExpr:
		return vm.evalUnaryExpr(e)
	case *ast.CallExpr:
		return vm.evalCallExpr(e)
	case *ast.MemberExpr:
		return vm.evalMemberExpr(e)
	case *ast.IndexExpr:
		return vm.evalIndexExpr(e)
	case *ast.TemplateLit:
		return vm.evalTemplateLit(e)
	default:
		return nil, NewRuntimeError(e.Pos(), "unknown expression type: %T", e)
	}
}

// Statement evaluation methods

func (vm *VM) evalLetStmt(stmt *ast.LetStmt) (Value, *RuntimeError) {
	value, err := vm.evalExpr(stmt.Value)
	if err != nil {
		return nil, err
	}
	vm.vars[stmt.Name.Name] = value
	return nil, nil
}

func (vm *VM) evalAssignStmt(stmt *ast.AssignStmt) (Value, *RuntimeError) {
	value, err := vm.evalExpr(stmt.Value)
	if err != nil {
		return nil, err
	}
	
	switch lhs := stmt.Target.(type) {
	case *ast.Ident:
		// Simple variable assignment
		if _, exists := vm.vars[lhs.Name]; !exists {
			return nil, NewRuntimeError(lhs.Pos(), "undefined variable: %s", lhs.Name)
		}
		vm.vars[lhs.Name] = value
		return value, nil
		
	case *ast.MemberExpr:
		// Object member assignment: obj.field = value
		obj, err := vm.evalExpr(lhs.Object)
		if err != nil {
			return nil, err
		}
		objMap, ok := obj.(Object)
		if !ok {
			return nil, NewRuntimeError(lhs.Pos(), "cannot assign to member of non-object")
		}
		objMap[lhs.Field.Name] = value
		return value, nil
		
	case *ast.IndexExpr:
		// Array/object index assignment: arr[i] = value or obj[key] = value
		target, err := vm.evalExpr(lhs.Target)
		if err != nil {
			return nil, err
		}
		index, err := vm.evalExpr(lhs.Index)
		if err != nil {
			return nil, err
		}
		
		if arr, ok := target.([]Value); ok {
			// Array assignment
			idx, ok := index.(float64)
			if !ok {
				return nil, NewRuntimeError(lhs.Pos(), "array index must be a number")
			}
			i := int(idx)
			if i < 0 || i >= len(arr) {
				return nil, NewRuntimeError(lhs.Pos(), "array index out of bounds: %d", i)
			}
			arr[i] = value
			return value, nil
		} else if obj, ok := target.(Object); ok {
			// Object assignment
			key, ok := index.(string)
			if !ok {
				return nil, NewRuntimeError(lhs.Pos(), "object key must be a string")
			}
			obj[key] = value
			return value, nil
		} else {
			return nil, NewRuntimeError(lhs.Pos(), "cannot index assign to non-array/non-object")
		}
		
	default:
		return nil, NewRuntimeError(stmt.Pos(), "invalid assignment target")
	}
}

func (vm *VM) evalExprStmt(stmt *ast.ExprStmt) (Value, *RuntimeError) {
	return vm.evalExpr(stmt.Value)
}

// Expression evaluation methods

func (vm *VM) evalNumberLit(lit *ast.NumberLit) (Value, *RuntimeError) {
	return lit.Value, nil
}

func (vm *VM) evalStringLit(lit *ast.StringLit) (Value, *RuntimeError) {
	return lit.Value, nil
}

func (vm *VM) evalIdent(ident *ast.Ident) (Value, *RuntimeError) {
	value, exists := vm.vars[ident.Name]
	if !exists {
		return nil, NewRuntimeError(ident.Pos(), "undefined variable: %s", ident.Name)
	}
	return value, nil
}

func (vm *VM) evalBinaryExpr(expr *ast.BinaryExpr) (Value, *RuntimeError) {
	left, err := vm.evalExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := vm.evalExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	
	switch expr.Operator {
	case "+":
		return vm.evalAdd(left, right, expr.Pos())
	case "-":
		return vm.evalSubtract(left, right, expr.Pos())
	case "*":
		return vm.evalMultiply(left, right, expr.Pos())
	case "/":
		return vm.evalDivide(left, right, expr.Pos())
	case "%":
		return vm.evalModulus(left, right, expr.Pos())
	case "==":
		return Equal(left, right), nil
	case "!=":
		return !Equal(left, right), nil
	case "<":
		return vm.evalLess(left, right, expr.Pos())
	case ">":
		return vm.evalGreater(left, right, expr.Pos())
	case "<=":
		return vm.evalLessEqual(left, right, expr.Pos())
	case ">=":
		return vm.evalGreaterEqual(left, right, expr.Pos())
	default:
		return nil, NewRuntimeError(expr.Pos(), "unknown binary operator: %s", expr.Operator)
	}
}

func (vm *VM) evalUnaryExpr(expr *ast.UnaryExpr) (Value, *RuntimeError) {
	operand, err := vm.evalExpr(expr.Operand)
	if err != nil {
		return nil, err
	}
	
	switch expr.Operator {
	case "-":
		if num, ok := operand.(float64); ok {
			return -num, nil
		}
		return nil, NewRuntimeError(expr.Pos(), "unary - requires a number")
	case "!":
		if b, ok := operand.(bool); ok {
			return !b, nil
		}
		return nil, NewRuntimeError(expr.Pos(), "unary ! requires a boolean")
	default:
		return nil, NewRuntimeError(expr.Pos(), "unknown unary operator: %s", expr.Operator)
	}
}

func (vm *VM) evalCallExpr(expr *ast.CallExpr) (Value, *RuntimeError) {
	callee, err := vm.evalExpr(expr.Callee)
	if err != nil {
		return nil, err
	}
	
	callable, ok := callee.(Callable)
	if !ok {
		return nil, NewRuntimeError(expr.Pos(), "value is not callable")
	}
	
	args := make([]Value, len(expr.Args))
	for i, argExpr := range expr.Args {
		arg, err := vm.evalExpr(argExpr)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}
	
	return callable.Call(expr.Pos(), args)
}

func (vm *VM) evalMemberExpr(expr *ast.MemberExpr) (Value, *RuntimeError) {
	obj, err := vm.evalExpr(expr.Object)
	if err != nil {
		return nil, err
	}
	
	if objMap, ok := obj.(Object); ok {
		value, exists := objMap[expr.Field.Name]
		if !exists {
			return nil, NewRuntimeError(expr.Pos(), "property '%s' not found", expr.Field.Name)
		}
		return value, nil
	}
	
	return nil, NewRuntimeError(expr.Pos(), "cannot access property of non-object")
}

func (vm *VM) evalIndexExpr(expr *ast.IndexExpr) (Value, *RuntimeError) {
	target, err := vm.evalExpr(expr.Target)
	if err != nil {
		return nil, err
	}
	index, err := vm.evalExpr(expr.Index)
	if err != nil {
		return nil, err
	}
	
	if arr, ok := target.([]Value); ok {
		// Array indexing
		idx, ok := index.(float64)
		if !ok {
			return nil, NewRuntimeError(expr.Pos(), "array index must be a number")
		}
		i := int(idx)
		if i < 0 || i >= len(arr) {
			return nil, NewRuntimeError(expr.Pos(), "array index out of bounds: %d", i)
		}
		return arr[i], nil
	} else if obj, ok := target.(Object); ok {
		// Object indexing
		key, ok := index.(string)
		if !ok {
			return nil, NewRuntimeError(expr.Pos(), "object key must be a string")
		}
		value, exists := obj[key]
		if !exists {
			return nil, NewRuntimeError(expr.Pos(), "key '%s' not found", key)
		}
		return value, nil
	}
	
	return nil, NewRuntimeError(expr.Pos(), "cannot index non-array/non-object")
}

func (vm *VM) evalTemplateLit(lit *ast.TemplateLit) (Value, *RuntimeError) {
	var result string
	for _, part := range lit.Parts {
		switch p := part.(type) {
		case *ast.TextPart:
			result += p.Value
		case *ast.Interpolation:
			value, err := vm.evalExpr(p.Expr)
			if err != nil {
				return nil, err
			}
			result += Stringify(value)
		}
	}
	return result, nil
}

// Binary operation helpers

func (vm *VM) evalAdd(left, right Value, pos domain.Pos) (Value, *RuntimeError) {
	if IsNumber(left) && IsNumber(right) {
		return left.(float64) + right.(float64), nil
	}
	if IsString(left) && IsString(right) {
		return left.(string) + right.(string), nil
	}
	return nil, NewRuntimeError(pos, "type mismatch for + operator")
}

func (vm *VM) evalSubtract(left, right Value, pos domain.Pos) (Value, *RuntimeError) {
	if IsNumber(left) && IsNumber(right) {
		return left.(float64) - right.(float64), nil
	}
	return nil, NewRuntimeError(pos, "- operator requires numbers")
}

func (vm *VM) evalMultiply(left, right Value, pos domain.Pos) (Value, *RuntimeError) {
	if IsNumber(left) && IsNumber(right) {
		return left.(float64) * right.(float64), nil
	}
	return nil, NewRuntimeError(pos, "* operator requires numbers")
}

func (vm *VM) evalDivide(left, right Value, pos domain.Pos) (Value, *RuntimeError) {
	if IsNumber(left) && IsNumber(right) {
		rightNum := right.(float64)
		if rightNum == 0 {
			return nil, NewRuntimeError(pos, "division by zero")
		}
		return left.(float64) / rightNum, nil
	}
	return nil, NewRuntimeError(pos, "/ operator requires numbers")
}

func (vm *VM) evalModulus(left, right Value, pos domain.Pos) (Value, *RuntimeError) {
	if IsNumber(left) && IsNumber(right) {
		rightNum := right.(float64)
		if rightNum == 0 {
			return nil, NewRuntimeError(pos, "modulus by zero")
		}
		leftNum := left.(float64)
		return float64(int64(leftNum) % int64(rightNum)), nil
	}
	return nil, NewRuntimeError(pos, "%% operator requires numbers")
}

func (vm *VM) evalLess(left, right Value, pos domain.Pos) (Value, *RuntimeError) {
	if IsNumber(left) && IsNumber(right) {
		return left.(float64) < right.(float64), nil
	}
	return nil, NewRuntimeError(pos, "< operator requires numbers")
}

func (vm *VM) evalGreater(left, right Value, pos domain.Pos) (Value, *RuntimeError) {
	if IsNumber(left) && IsNumber(right) {
		return left.(float64) > right.(float64), nil
	}
	return nil, NewRuntimeError(pos, "> operator requires numbers")
}

func (vm *VM) evalLessEqual(left, right Value, pos domain.Pos) (Value, *RuntimeError) {
	if IsNumber(left) && IsNumber(right) {
		return left.(float64) <= right.(float64), nil
	}
	return nil, NewRuntimeError(pos, "<= operator requires numbers")
}

func (vm *VM) evalGreaterEqual(left, right Value, pos domain.Pos) (Value, *RuntimeError) {
	if IsNumber(left) && IsNumber(right) {
		return left.(float64) >= right.(float64), nil
	}
	return nil, NewRuntimeError(pos, ">= operator requires numbers")
}

// Default implementations for load/save (can be overridden)
func (vm *VM) defaultLoad(path string) (*Map, error) {
	return nil, fmt.Errorf("load not implemented")
}

func (vm *VM) defaultSave(m *Map, path string) error {
	return fmt.Errorf("save not implemented")
}
