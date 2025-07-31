# WJS Virtual Machine Specification

**Language**: WJS (Worldographer JavaScript-ish)  

**Purpose**: Defines the runtime execution model for WJS. This document describes how the VM should evaluate a parsed AST and manage program state.

## Overview

The VM interprets a parsed `Program` (AST) and executes statements and expressions.
It supports basic arithmetic, variables, object access, indexing, function calls, and template string interpolation.
The runtime is single-pass and fail-fast: any error immediately halts execution.

The VM provides access to built-in functions and external data structures representing Worldographer maps.

---

## Execution Model

- The VM walks the AST in evaluation order.
- The VM maintains a global environment mapping identifiers to values.
- Values are boxed as `Value` types with support for runtime type inspection.
- All execution errors must be reported with position metadata from the AST.
- Evaluation is non-concurrent and non-threaded.

---

## Types

The VM must support the following value types:

| Go Type   | WJS Type  |
|-----------|-----------|
| `int64`   | `number`  |
| `float64` | `number`  |
| `string`  | `string`  |
| `bool`    | `boolean` |
| `nil`     | `null`    |
| `error`   | `error`   |
| `[]Value` | `array`   |
| `*Object` | `object`  |
| `*Map`    | `Map`     |

Use a common boxed interface:

```go
type Value interface{}

type Object map[string]Value
````

The `Map` type is a Go struct loaded from `.wxx` files, passed through `load()` and `save()` built-ins.

---

## Environment

The environment is a flat key-value map (no block scoping):

```go
type Env struct {
    vars map[string]Value
}
```

* `let x = ...` binds a value
* `x = ...` updates an existing binding or raises an error if not declared

---

## Statements

### LetStmt

Evaluate `Value` and bind to `Name` in the environment.

### AssignStmt

Evaluate LHS as a reference (ident, member, or index). Assign evaluated RHS to that location.

### ExprStmt

Evaluate the expression for side effects. Result is ignored.

---

## Expressions

### BinaryExpr

| Operator             | Behavior                     |
| -------------------- | ---------------------------- |
| `+`                  | numeric add or string concat |
| `-`                  | numeric subtraction          |
| `*`                  | numeric multiply             |
| `/`                  | numeric divide               |
| `%`                  | numeric modulus              |
| `==`                 | deep equality                |
| `!=`                 | deep inequality              |
| `<`, `>`, `<=`, `>=` | numeric comparison           |

Operands must be of the same type. Mixed types are runtime errors.

### UnaryExpr

| Operator | Behavior         |
| -------- | ---------------- |
| `-`      | numeric negation |
| `!`      | boolean negation |

### CallExpr

* Evaluate callee: must be a function (`Value` satisfying a `Callable` interface).
* Evaluate arguments.
* Call the function with evaluated arguments.
* If the function panics or fails, halt execution and return a VM error.

Built-in functions:

```go
func load(path string) (*Map, error)
func save(map *Map, path string) error
func print(args ...Value)
```

VM provides these under identifiers `load`, `save`, and `print`.

### MemberExpr

* Evaluate `Object`
* Look up field using identifier name
* Object must be `map[string]Value`, `*Object`, or known struct with reflect access (e.g., `*Map`)
* If key not found, raise runtime error

### IndexExpr

* Evaluate target and index
* Target must be an array or object
* Index must be numeric (array) or string (object)
* Out-of-bounds or missing keys result in a runtime error

### TemplateLit

* Evaluate all parts
* Concatenate into a single string
* Interpolation parts must be converted to strings via `fmt.Sprintf("%v", ...)`

---

## Errors

All errors should be represented using:

```go
type RuntimeError struct {
    Pos     Pos
    Message string
}
```

The VM must immediately halt on the first error and return it to the caller.
No panics should escape evaluation.

---

## Entry Point

The VM should expose a single entry function:

```go
func Eval(prog *ast.Program, env *Env) (Value, *RuntimeError)
```

* Evaluates the program
* Returns last expression result (if any) and a runtime error (if any)

---

## Testing

Unit tests should include:

* Arithmetic and logic expressions
* Assignments to identifiers, object fields, and array indices
* Call and side effects of built-ins (`print`, `load`, `save`)
* Failure modes: missing fields, invalid types, undefined identifiers
* Template strings with and without expressions

Test helpers should allow capturing output from `print(...)`.

---

## Out of Scope

The following are not required for v1:

* User-defined functions or lambdas
* Closures or nested scopes
* Garbage collection
* Control flow (`if`, `while`, etc.)
* Importing external code
* Reflection on user-defined objects

---

End of specification.

```
