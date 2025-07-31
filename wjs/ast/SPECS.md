# WJS AST Specification

**Language:** WJS (Worldographer JavaScript-ish)  
**Purpose:** This document specifies the Abstract Syntax Tree (AST) node types for the WJS scripting language. It is intended for Go developers implementing the parser and interpreter.

---

## 📘 Overview

The AST represents parsed WJS source code as a tree of nodes. Each node implements the base `Node` interface and tracks its source position for error reporting.

```go
type Node interface {
    Pos() domain.Pos
}
```

---

## 🏗️ Base Interfaces

### Node Categories

```go
type Stmt interface {
    Node
    isStmt()
}

type Expr interface {
    Node
    isExpr()
}
```

All statements implement `Stmt`, all expressions implement `Expr`. This enables type-safe pattern matching in the parser and interpreter.

---

## 📄 Statement Nodes

### Program (Root Node)

```go
type Program struct {
    Start domain.Pos
    Stmts []Stmt
}
```

* Root of every AST
* Contains zero or more statements
* `Start` position is typically beginning of file

### Let Statement

```go
type LetStmt struct {
    Start domain.Pos
    Name  *Ident
    Value Expr
}
```

* Variable declaration with initialization
* Example: `let map = load("world.wxx");`
* `Name` is the variable identifier
* `Value` is the initializing expression

### Assignment Statement

```go
type AssignStmt struct {
    Start  domain.Pos
    Target Expr // must be Ident, IndexExpr, or MemberExpr
    Value  Expr
}
```

* Assigns value to existing variable, property, or array element
* Example: `map.tiles[2][4] = "mountain";`
* `Target` must be a valid left-hand side (validated by `CheckValid`)

### Expression Statement

```go
type ExprStmt struct {
    Start domain.Pos
    Value Expr
}
```

* Expression used as a statement
* Example: `print("Hello world");`
* Common for function calls with side effects

---

## 🧮 Expression Nodes

### Literals

#### Identifier

```go
type Ident struct {
    Start domain.Pos
    Name  string
}
```

* Variable or function name
* Example: `map`, `load`, `width`

#### Number Literal

```go
type NumberLit struct {
    Start domain.Pos
    Value float64
}
```

* Numeric constant (integer or floating-point)
* Examples: `42`, `3.14`, `-0.001`

#### String Literal

```go
type StringLit struct {
    Start domain.Pos
    Value string
}
```

* String constant with quotes removed
* Examples: `"world.wxx"`, `'mountain'`

#### Template Literal

```go
type TemplateLit struct {
    Start domain.Pos
    Parts []TemplatePart
}
```

* Template string with interpolation
* Example: `` `Map has ${map.width} tiles` ``
* Parts alternate between text and interpolated expressions

### Template Parts

```go
type TemplatePart interface {
    Node
    isTemplatePart()
}
```

#### Text Part

```go
type TextPart struct {
    Start domain.Pos
    Value string
}
```

* Raw text within template literal
* Example: `"Map has "` in `` `Map has ${width}` ``

#### Interpolation

```go
type Interpolation struct {
    Start domain.Pos
    Expr  Expr
}
```

* Expression within `${...}` in template literal
* Example: `map.width` in `` `Map has ${map.width}` ``

---

## 🛠️ Composite Expressions

### Binary Expression

```go
type BinaryExpr struct {
    Start    domain.Pos
    Left     Expr
    Operator string // "+", "-", "*", "/", "%", "==", "!=", "<", ">", "<=", ">="
    Right    Expr
}
```

* Two operands with infix operator
* Examples: `x + y`, `map.width * map.height`, `value == null`

### Unary Expression

```go
type UnaryExpr struct {
    Start    domain.Pos
    Operator string // "-" or "!"
    Operand  Expr
}
```

* Single operand with prefix operator
* Examples: `-value`, `!found`

### Function Call

```go
type CallExpr struct {
    Start  domain.Pos
    Callee Expr // usually Ident
    Args   []Expr
}
```

* Function invocation
* Examples: `load("world.wxx")`, `print(message, value)`
* `Callee` is typically an identifier, but could be any expression

### Member Access

```go
type MemberExpr struct {
    Start  domain.Pos
    Object Expr
    Field  *Ident
}
```

* Property access using dot notation
* Examples: `map.width`, `tile.terrain`
* `Field` must be a valid identifier

### Index Access

```go
type IndexExpr struct {
    Start  domain.Pos
    Target Expr
    Index  Expr
}
```

* Array/object element access using bracket notation
* Examples: `tiles[x]`, `map.tiles[x][y]`
* `Index` can be any expression that evaluates to a valid key

---

## 🧪 AST Validation

The `CheckValid` function performs semantic validation:

### Valid Assignment Targets

Only these expressions can appear on the left side of assignments:
* `Ident` - simple variable
* `MemberExpr` - property access
* `IndexExpr` - array/object element

### Validation Rules

* Identifiers must have non-empty names
* Binary expressions must have both operands
* Unary expressions must have an operand
* Template literals must have at least one part
* Interpolations must contain an expression
* Member expressions must have valid field names

---

## 🎯 Usage Examples

### Simple Variable Declaration

```js
let x = 42;
```

AST Structure:
```
Program
  LetStmt x =
    Number 42
```

### Function Call with Template String

```js
print(`Value is ${x + 1}`);
```

AST Structure:
```
Program
  ExprStmt
    CallExpr
      Ident "print"
        Template
          Text "Value is "
          Interpolation
            BinaryExpr "+"
              Ident "x"
              Number 1
```

### Property Assignment

```js
map.tiles[x][y].terrain = "mountain";
```

AST Structure:
```
Program
  AssignStmt
    MemberExpr
      IndexExpr
        IndexExpr
          MemberExpr
            Ident "map"
            Ident "tiles"
          Ident "x"
        Ident "y"
      Ident "terrain"
    String "mountain"
```

---

## 🧰 Helper Functions

### Pretty Printing

Use `DumpAST(node)` or `PrettyPrint(node)` to visualize the AST structure for debugging.

### Position Tracking

All nodes track their source position via `domain.Pos`:

```go
type Pos struct {
    Line   int    // 1-based
    Column int    // 1-based
    Offset int    // byte offset in input
    Path   string // source file name
}
```

This enables precise error reporting during parsing and interpretation.

---

## 📂 Implementation Notes

* The AST is designed for single-pass interpretation
* All nodes are immutable after construction
* Position information is mandatory for error reporting
* The type system prevents invalid AST construction through Go's type checker
* Template literal parsing happens in the lexer; the AST just stores the structured parts

---

## 🛑 Limitations

* No function declarations or blocks (flat scope only)
* No control flow statements (`if`, `while`, `for`)
* No object/array literal expressions
* No destructuring or advanced assignment patterns

These limitations are intentional for the MVP and may be lifted in future versions.
