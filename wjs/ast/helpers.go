// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ast

import (
	"fmt"
	"strings"
)

func DumpAST(n Node) {
	fmt.Println(PrettyPrint(n))
}

// todo: consider putting position into the pretty print like this:
//   fmt.Fprintf(b, "%sLetStmt (%d:%d) %s =\n", indentStr, node.Start.Line, node.Start.Column, node.Name.Name)

func PrettyPrint(n Node) string {
	var b strings.Builder
	writePretty(&b, n, 0)
	return b.String()
}

func writePretty(b *strings.Builder, n Node, indent int) {
	indentStr := strings.Repeat("  ", indent)

	switch node := n.(type) {
	case *Program:
		b.WriteString("Program\n")
		for _, stmt := range node.Stmts {
			writePretty(b, stmt, indent+1)
		}

	case *LetStmt:
		fmt.Fprintf(b, "%sLetStmt %s =\n", indentStr, node.Name.Name)
		writePretty(b, node.Value, indent+1)

	case *AssignStmt:
		fmt.Fprintf(b, "%sAssignStmt\n", indentStr)
		writePretty(b, node.Target, indent+1)
		writePretty(b, node.Value, indent+1)

	case *ExprStmt:
		fmt.Fprintf(b, "%sExprStmt\n", indentStr)
		writePretty(b, node.Value, indent+1)

	case *Ident:
		fmt.Fprintf(b, "%sIdent %q\n", indentStr, node.Name)

	case *NumberLit:
		if node.IntVal != nil {
			fmt.Fprintf(b, "%sNumber %d\n", indentStr, *node.IntVal)
		} else if node.FloatVal != nil {
			fmt.Fprintf(b, "%sNumber %v\n", indentStr, *node.FloatVal)
		} else {
			fmt.Fprintf(b, "%sNumber <invalid>\n", indentStr)
		}

	case *StringLit:
		fmt.Fprintf(b, "%sString %q\n", indentStr, node.Value)

	case *TemplateLit:
		fmt.Fprintf(b, "%sTemplate\n", indentStr)
		for _, part := range node.Parts {
			writePretty(b, part, indent+1)
		}

	case *TextPart:
		fmt.Fprintf(b, "%sText %q\n", indentStr, node.Value)

	case *Interpolation:
		fmt.Fprintf(b, "%sInterpolation\n", indentStr)
		writePretty(b, node.Expr, indent+1)

	case *BinaryExpr:
		fmt.Fprintf(b, "%sBinaryExpr %q\n", indentStr, node.Operator)
		writePretty(b, node.Left, indent+1)
		writePretty(b, node.Right, indent+1)

	case *UnaryExpr:
		fmt.Fprintf(b, "%sUnaryExpr %q\n", indentStr, node.Operator)
		writePretty(b, node.Operand, indent+1)

	case *CallExpr:
		fmt.Fprintf(b, "%sCallExpr\n", indentStr)
		writePretty(b, node.Callee, indent+1)
		for _, arg := range node.Args {
			writePretty(b, arg, indent+2)
		}

	case *MemberExpr:
		fmt.Fprintf(b, "%sMemberExpr\n", indentStr)
		writePretty(b, node.Object, indent+1)
		writePretty(b, node.Field, indent+1)

	case *IndexExpr:
		fmt.Fprintf(b, "%sIndexExpr\n", indentStr)
		writePretty(b, node.Target, indent+1)
		writePretty(b, node.Index, indent+1)

	default:
		fmt.Fprintf(b, "%s<unknown node type>\n", indentStr)
	}
}

// CheckValid walks the AST and returns the first semantic error found, or nil if valid.
//
// * Ensures LHS of assignments is valid (let x = 1 ✅, 1 = x ❌)
//
// * Validates Ident names are not empty
//
// * Enforces template strings are non-empty and interpolation contains valid expressions
//
// * Recursively checks expression subtrees
func CheckValid(n Node) error {
	switch node := n.(type) {
	case *Program:
		for _, stmt := range node.Stmts {
			if err := CheckValid(stmt); err != nil {
				return err
			}
		}

	case *LetStmt:
		if node.Name == nil || node.Name.Name == "" {
			return fmt.Errorf("invalid let statement at %d:%d: missing variable name", node.Start.Line, node.Start.Column)
		}
		return CheckValid(node.Value)

	case *AssignStmt:
		if err := checkValidLHS(node.Target); err != nil {
			return fmt.Errorf("invalid assignment target at %d:%d: %w", node.Start.Line, node.Start.Column, err)
		}
		return CheckValid(node.Value)

	case *ExprStmt:
		return CheckValid(node.Value)

	case *BinaryExpr:
		if node.Left == nil || node.Right == nil {
			return fmt.Errorf("incomplete binary expression at %d:%d", node.Start.Line, node.Start.Column)
		}
		if err := CheckValid(node.Left); err != nil {
			return err
		}
		if err := CheckValid(node.Right); err != nil {
			return err
		}

	case *UnaryExpr:
		if node.Operand == nil {
			return fmt.Errorf("missing operand in unary expression at %d:%d", node.Start.Line, node.Start.Column)
		}
		return CheckValid(node.Operand)

	case *CallExpr:
		if err := CheckValid(node.Callee); err != nil {
			return err
		}
		for _, arg := range node.Args {
			if err := CheckValid(arg); err != nil {
				return err
			}
		}

	case *MemberExpr:
		if err := CheckValid(node.Object); err != nil {
			return err
		}
		if node.Field == nil || node.Field.Name == "" {
			return fmt.Errorf("invalid member field at %d:%d", node.Start.Line, node.Start.Column)
		}

	case *IndexExpr:
		if err := CheckValid(node.Target); err != nil {
			return err
		}
		if err := CheckValid(node.Index); err != nil {
			return err
		}

	case *TemplateLit:
		if len(node.Parts) == 0 {
			return fmt.Errorf("empty template string at %d:%d", node.Start.Line, node.Start.Column)
		}
		for _, part := range node.Parts {
			if err := CheckValid(part); err != nil {
				return err
			}
		}

	case *Interpolation:
		if node.Expr == nil {
			return fmt.Errorf("missing expression in interpolation at %d:%d", node.Start.Line, node.Start.Column)
		}
		return CheckValid(node.Expr)

	case *TextPart:
		// No validation needed.

	case *Ident:
		if node.Name == "" {
			return fmt.Errorf("empty identifier at %d:%d", node.Start.Line, node.Start.Column)
		}

	case *NumberLit, *StringLit:
		// Always valid.

	default:
		return fmt.Errorf("unknown or unsupported AST node at %v", n.Pos())
	}

	return nil
}

func checkValidLHS(e Expr) error {
	switch e.(type) {
	case *Ident, *MemberExpr, *IndexExpr:
		return nil
	default:
		return fmt.Errorf("invalid left-hand side: must be identifier, member, or index")
	}
}
