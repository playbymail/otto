// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ast

type Pos struct {
	Line   int
	Column int
}

type Program struct {
	Statements []Stmt
	Pos        Pos
}

type Stmt interface {
	Position() Pos
}

type Statement struct {
	Pos Pos
}

type AssignStmt struct {
	Target Expr
	Value  Expr
	Pos    Pos
}

func (s *AssignStmt) Position() Pos { return s.Pos }

type ExprStmt struct {
	Expr Expr
	Pos  Pos
}

func (s *ExprStmt) Position() Pos { return s.Pos }

type LetStmt struct {
	Name  *Ident
	Value Expr
	Pos   Pos
}

func (s *LetStmt) Position() Pos { return s.Pos }

type Expr interface {
	Position() Pos // or any method that lets you identify this as an Expr
}

type BinaryExpr struct {
	Left     Expr
	Operator string
	Right    Expr
	Pos      Pos
}

func (b *BinaryExpr) Position() Pos { return b.Pos }

type CallExpr struct {
	Callee Expr // usually an *Ident, but could be more complex in full implementation
	Args   []Expr
	Pos    Pos
}

func (e *CallExpr) Position() Pos { return e.Pos }

type IndexExpr struct {
	Target Expr
	Index  Expr
	Pos    Pos
}

func (e *IndexExpr) Position() Pos { return e.Pos }

type MemberExpr struct {
	Object Expr
	Name   *Ident
	Pos    Pos
}

func (e *MemberExpr) Position() Pos { return e.Pos }

type UnaryExpr struct {
	Operator string
	Expr     Expr
	Pos      Pos
}

func (u *UnaryExpr) Position() Pos { return u.Pos }

// Suffix represents a postfix operation applied to an expression.
//
// It is used to model chained accessors such as function calls,
// member accesses, and index operations in expressions like:
//
//	foo.bar(1)[2]
//
// Each Suffix implementation applies itself to a base expression
// and returns a new expression node representing the combined form.
//
// For example:
//   - CallSuffix applies function call arguments to a base.
//   - IndexSuffix applies a bracketed index to a base.
//   - MemberSuffix applies a dot-accessed field to a base.
type Suffix interface {
	// Apply transforms the base expression by applying the suffix,
	// returning a new expression node.
	Apply(base Expr) Expr
}

type CallSuffix struct {
	Args []Expr
	Pos  Pos
}

func (c *CallSuffix) Apply(base Expr) Expr {
	return &CallExpr{
		Callee: base,
		Args:   c.Args,
		Pos:    c.Pos,
	}
}

type IndexSuffix struct {
	Index Expr
	Pos   Pos
}

func (i *IndexSuffix) Apply(base Expr) Expr {
	return &IndexExpr{
		Target: base,
		Index:  i.Index,
		Pos:    i.Pos,
	}
}

type MemberSuffix struct {
	Object Expr // the base expression
	Name   *Ident
	Pos    Pos
}

func (m *MemberSuffix) Apply(base Expr) Expr {
	return &MemberExpr{
		Object: base,
		Name:   m.Name,
		Pos:    m.Pos,
	}
}

// literals

type BoolLiteral struct {
	Value bool
	Pos   Pos
}

func (b *BoolLiteral) Position() Pos { return b.Pos }

type NullLiteral struct {
	Pos Pos
}

func (n *NullLiteral) Position() Pos { return n.Pos }

type NumberLiteral struct {
	Value any // user must cast to float64 or int64!
	// consider adding `IsFloat bool` to help casting
	Pos Pos
}

func (n *NumberLiteral) Position() Pos { return n.Pos }

type StringLiteral struct {
	Value string
	Pos   Pos
}

func (s *StringLiteral) Position() Pos { return s.Pos }

type Ident struct {
	Name string
	Pos  Pos
}

func (i *Ident) Position() Pos { return i.Pos }
