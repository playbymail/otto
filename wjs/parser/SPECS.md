# WJS Parser Specification

**Language**: WJS (Worldographer JavaScript-ish)  

**Purpose**: Defines the grammar and structural rules for WJS source code.
This document is for developers implementing the parser using a recursive descent strategy in Go.

## Overview

The parser consumes a token stream produced by the WJS lexer and produces an abstract syntax tree (AST) using the types defined in `wjs/ast`.

The grammar is designed to be LL(1) friendly and avoids left recursion.
It uses semicolons as explicit statement terminators.
The grammar includes variable declarations, assignments, expression statements, property/index access, function calls, string interpolation, and basic literals.

The parser must preserve position information (`Pos`) from tokens and attach it to all AST nodes for accurate error reporting.

---

## Source Units

A WJS file is parsed as a single `Program`.

```grammar

Program     = { Statement } ;

```

---

## Statements

All statements must terminate with a semicolon.

```grammar

Statement   = LetStmt
            | AssignStmt
            | ExprStmt ;

LetStmt     = "let" Identifier "=" Expression ";" ;

AssignStmt  = Target "=" Expression ";" ;

ExprStmt    = Expression ";" ;

```

---

## Assignment Targets

Targets may be identifiers, property accesses, or index expressions.

```grammar

Target      = Identifier [ "." Identifier | "[" Expression "]" ] ;

```

The parser must verify that the left-hand side of an assignment is a valid assignable expression.
These include:

- `Ident`
- `MemberExpr`
- `IndexExpr`

---

## Expressions

Expressions support arithmetic, comparisons, property access, indexing, function calls, literals, and interpolated strings.

```grammar

Expression  = Equality ;

Equality    = Comparison [ ( "==" | "!=" ) Comparison ] ;

Comparison  = Term [ ( "<" | ">" | "<=" | ">=" ) Term ] ;

Term        = Factor { ( "+" | "-" ) Factor } ;

Factor      = Unary { ( "\*" | "/" | "%" ) Unary } ;

Unary       = ( "-" | "!" )? Primary ;

```

---

## Primary Expressions

Primary expressions include literals, identifiers, parenthesized expressions, and composite forms like function calls, member access, and indexing.

```grammar

Primary     = Literal
            | Identifier
            | "(" Expression ")"
            | CallExpr
            | IndexExpr
            | MemberExpr ;

CallExpr    = Primary "(" [ ExpressionList ] ")" ;

IndexExpr   = Primary "[" Expression "]" ;

MemberExpr  = Primary "." Identifier ;

ExpressionList = Expression { "," Expression } ;

```

Function calls, indexing, and member expressions are all postfix operations on primary expressions and must be parsed iteratively:

Example:

```

map.tiles[0][1].terrain

```

Is parsed as:

```

MemberExpr(
  IndexExpr(
    IndexExpr(
      MemberExpr(
        Ident("map"),
        Ident("tiles")
      ),
      NumberLit(0)
    ),
    NumberLit(1)
  ),
  Ident("terrain")
)

```

---

## Literals

```grammar

Literal     = Number | String | Template | "true" | "false" | "null" ;

Identifier  = letter { letter | digit | "\_" } ;

Number      = digit { digit } [ "." digit { digit } ] ;

String      = '"' { any character except '"' } '"'
            | "'" { any character except "'" } "'" ;

Template    = "`" { TemplatePart | "${" Expression "}" } "`" ;

TemplatePart = any character except "\`" and "\${" ;

````

Template strings are parsed as a single literal token by the lexer.
The parser is responsible for splitting the template into a sequence of text and embedded expressions.
Each template is lowered to a `TemplateLit` node containing `[]TemplatePart`, where `TemplatePart` is either a `TextPart` or `Interpolation`.

---

## Error Handling

The parser must fail fast and return immediately on the first error encountered.
Each error must include:

- Line and column from the token or node that caused the error
- A brief diagnostic message
- Optionally, a snippet of source context (helpful in REPL or tool integration)

The parser must not panic.
Errors should be represented using a structured error type that includes position metadata.

---

## AST Construction

The parser must produce an abstract syntax tree conforming to the `ast` package.
All nodes must include position information using the `Pos` struct:

```go
type Pos struct {
    Line   int
    Column int
    Offset int
}
````

Every AST node (statement or expression) must implement the `Node` interface:

```go
type Node interface {
    Pos() Pos
}
```

Statement and expression nodes must also implement `Stmt` or `Expr`.

---

## Semicolon Rules

Semicolons are **required** after all statements.
There is no automatic semicolon insertion.
This simplifies parsing and eliminates ambiguity.

---

## Disallowed Constructs (Out of Scope)

The following features are explicitly not part of this version of the grammar:

* Blocks or scopes (no `{}` blocks)
* Function declarations
* Loops or control flow beyond `if`/`else`
* Arrow functions or anonymous functions
* Compound assignments (`+=`, `--`, etc.)
* Ternary expressions
* Chained comparison (`a < b < c`)
* Optional chaining or nullish coalescing

---

## Implementation Notes

* Parse expressions using precedence climbing or recursive methods as appropriate.
* Keep the parser deterministic and LL(1). Backtracking is not wanted.
* Use lookahead tokens to disambiguate between primary/postfix expressions (e.g. `(` for call, `[` for index, `.` for member).
  * For example, the parser must **disambiguate Primary → CallExpr vs MemberExpr vs IndexExpr** by peeking ahead after `Primary`.
* Template literal parsing should include a post-tokenization step to divide the template string into static text and `${}` interpolations.
* All statements must end in a semicolon (`;`), including expression statements.
* The `Template` rule is parsed as a single literal by the lexer and split into parts by a post-pass (not recursive descent).
* We assume a flat scope for now — no blocks or function declarations.
* No need yet for `return`, `while`, or `function` - this is an expression-oriented scripting language.

---

## Example

Given this input:

```js
let map = load("world.wxx");
map.tiles[0][1].terrain = "mountain";
print(`Done updating ${map.name}.`);
```

The AST should include:

* A `LetStmt` assigning the result of a `CallExpr` to `map`
* An `AssignStmt` updating a nested `MemberExpr` of an `IndexExpr`
* An `ExprStmt` for the `print` call, using a `TemplateLit`

---
