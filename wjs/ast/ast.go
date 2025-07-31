// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ast

import (
	"github.com/playbymail/otto/wjs/domain"
)

// 🧱 Base Node Interfaces

type Node interface {
	Pos() domain.Pos
}

type Stmt interface {
	Node
	isStmt()
}

type Expr interface {
	Node
	isExpr()
}

// 📄 Statement Nodes

type LetStmt struct {
	Start domain.Pos
	Name  *Ident
	Value Expr
}

func (s *LetStmt) Pos() domain.Pos { return s.Start }
func (s *LetStmt) isStmt()         {}

type AssignStmt struct {
	Start  domain.Pos
	Target Expr // must be Ident, IndexExpr, or MemberExpr
	Value  Expr
}

func (s *AssignStmt) Pos() domain.Pos { return s.Start }
func (s *AssignStmt) isStmt()         {}

type ExprStmt struct {
	Start domain.Pos
	Value Expr
}

func (s *ExprStmt) Pos() domain.Pos { return s.Start }
func (s *ExprStmt) isStmt()         {}

// 🧮 Expression Nodes

type Ident struct {
	Start domain.Pos
	Name  string
}

func (e *Ident) Pos() domain.Pos { return e.Start }
func (e *Ident) isExpr()         {}

type NumberLit struct {
	Start domain.Pos
	Value float64
}

func (e *NumberLit) Pos() domain.Pos { return e.Start }
func (e *NumberLit) isExpr()         {}

type StringLit struct {
	Start domain.Pos
	Value string
}

func (e *StringLit) Pos() domain.Pos { return e.Start }
func (e *StringLit) isExpr()         {}

type TemplateLit struct {
	Start domain.Pos
	Parts []TemplatePart // e.g., ["foo", expr, "bar"]
}

func (e *TemplateLit) Pos() domain.Pos { return e.Start }
func (e *TemplateLit) isExpr()         {}

type TemplatePart interface {
	Node
	isTemplatePart()
}

type TextPart struct {
	Start domain.Pos
	Value string
}

func (p *TextPart) Pos() domain.Pos { return p.Start }
func (p *TextPart) isTemplatePart() {}

type Interpolation struct {
	Start domain.Pos
	Expr  Expr
}

func (p *Interpolation) Pos() domain.Pos { return p.Start }
func (p *Interpolation) isTemplatePart() {}

// 🛠️ Composite Expressions

type BinaryExpr struct {
	Start    domain.Pos
	Left     Expr
	Operator string // "+", "-", "==", etc.
	Right    Expr
}

func (e *BinaryExpr) Pos() domain.Pos { return e.Start }
func (e *BinaryExpr) isExpr()         {}

type UnaryExpr struct {
	Start    domain.Pos
	Operator string // "-" or "!"
	Operand  Expr
}

func (e *UnaryExpr) Pos() domain.Pos { return e.Start }
func (e *UnaryExpr) isExpr()         {}

type CallExpr struct {
	Start  domain.Pos
	Callee Expr // usually Ident
	Args   []Expr
}

func (e *CallExpr) Pos() domain.Pos { return e.Start }
func (e *CallExpr) isExpr()         {}

type MemberExpr struct {
	Start  domain.Pos
	Object Expr
	Field  *Ident
}

func (e *MemberExpr) Pos() domain.Pos { return e.Start }
func (e *MemberExpr) isExpr()         {}

type IndexExpr struct {
	Start  domain.Pos
	Target Expr
	Index  Expr
}

func (e *IndexExpr) Pos() domain.Pos { return e.Start }
func (e *IndexExpr) isExpr()         {}

// 📦 Root Node

type Program struct {
	Start domain.Pos
	Stmts []Stmt
}

func (p *Program) Pos() domain.Pos { return p.Start }
