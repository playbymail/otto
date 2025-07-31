# WJS
WJS is a script for working with Worldographer files.

## ✅ **Initial Language Goals for `wjs`**

### 🎯 Purpose

* Read and modify Worldographer `.wxx` files.
* Enable basic scripting operations over map data.
* Use a syntax similar to Javascript to decrease learning time for new users.

---

## 🧱 Core Features

### 1. **Variable Declarations**

```js
let map = load("map.wxx")
let name = map.name
let area = map.width * map.height
```

* Only `let` is supported.
* No `var`, `const`, or hoisting.
* RHS expressions: literals, object/array literals, binary ops, function calls, property access.

---

### 2. **Property and Index Access**

```js
let terrain = map.tiles[5][3].terrain
map.tiles[2][4].terrain = "mountain"
```

* Dot notation (`obj.key`)
* Bracket notation (`arr[i]`)
* Nested indexing allowed

---

### 3. **Collections**

```js
let center = map.tiles[map.width / 2][map.height / 2]
let neighbors = [north, south, east, west]
```

* Array literals: `[a, b, c]`
* Object literals: `{ key: value }` (optional at first)
* Optional: simple `for` over arrays later

---

### 4. **Print Function with Interpolation**

```js
print(`Center is at ${map.width/2}, ${map.height/2}`)
```

* Interpolated strings (`backticks`, `${expr}` inside)
* Only one built-in function at first: `print(...)`

---

### 5. **Function Calls**

```js
let map = load("world.wxx")
save(map, "out.wxx")
```

* Built-in functions:

    * `load(path: string): Map`
    * `save(map: Map, path: string): void`
    * `print(...args: any[]): void`

---

## 🎭 Example Script

```js
let map = load("world.wxx")
let tile = map.tiles[10][5]
print(`Tile terrain is ${tile.terrain}`)

tile.terrain = "mountain"
save(map, "mountains.wxx")
```

---

## 🧠 Next Steps

1. **Grammar Design** (subset of JS):

    * Start with expressions: identifiers, literals, indexing, property access, function calls
    * Then add statements: `let`, assignment, expression statements

2. **AST Nodes**:

    * `Program`, `LetStmt`, `AssignStmt`, `ExprStmt`
    * `Ident`, `NumberLit`, `StringLit`, `BinaryExpr`, `CallExpr`, `IndexExpr`, `MemberExpr`, `InterpolatedString`

3. **Parser**:

    * Recursive descent: simple precedence climbing for expressions
    * Track source positions for error reporting

4. **Interpreter/VM**:

    * Simple environment map for `let` bindings
    * Built-ins for `load`, `save`, `print`
    * Return early on errors with source position

---

Would you like to start with the grammar and token list first, or jump into the AST scaffolding for `wjs`?

---

## ✅ 1. **Simplest Error Handling (Fail-Fast)**

> "If any error occurs, abort execution and print a meaningful message."

### Behavior:

* Any error in `load`, `save`, or invalid operations immediately halts the script.
* Error includes: line/column, expression, and message.

### Example Output:

```
Error at line 2, column 9: load("missing.wxx")
→ File not found: missing.wxx
```

### Pros:

* Very simple to implement.
* Encourages short scripts that succeed or fail cleanly.
* Ideal for initial REPL or batch mode.

### Cons:

* No recovery from errors.
* Script writers can't handle errors themselves.

**Recommendation:** Use this as your **default mode** to get up and running.

---

Would you like a proposal for the `RuntimeError` structure and how we might report it consistently across the VM and REPL?

---

