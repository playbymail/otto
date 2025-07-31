# WJS Lexer Specification

**Language:** WJS (Worldographer JavaScript-ish)  
**Purpose:** This document specifies the lexical rules for the WJS scripting language. It is intended for Go developers implementing the lexer.

---

## 📘 Overview

The lexer reads UTF-8 source code and emits a stream of tokens, each annotated with a `Pos` struct:

```go
type Token struct {
    Type    TokenType
    Lexeme  string
    Pos     Pos
}
````

```go
type Pos struct {
    Line   int // 1-based
    Column int // 1-based
    Offset int // byte offset in input
    Path   string // name of the source file (set only when processing a script file)
}
```

---

## 🎯 Goals

* Lexical simplicity for recursive descent parsing
* JavaScript-like feel with a *strict* and *minimal* token set
* Clear and consistent error handling for illegal characters and unterminated constructs

---

## 🎫 Token Types

### Punctuation

| Lexeme | TokenType       |
| ------ |-----------------|
| `(`    | `LPAREN`        |
| `)`    | `RPAREN`        |
| `{`    | `LBRACE`        |
| `}`    | `RBRACE`        |
| `[`    | `LBRACK`        |
| `]`    | `RBRACK`        |
| `.`    | `DOT`           |
| `,`    | `COMMA`         |
| `;`    | `SEMICOLON`     |
| `=`    | `EQUAL`         |
| `:`    | `COLON`         |
| `` ` `` | `BACKTICK`      |
| `${` |  `DOLLARBRACK` |


---

### Operators

| Lexeme | TokenType  |
| ------ | ---------- |
| `+`    | `PLUS`     |
| `-`    | `MINUS`    |
| `*`    | `ASTERISK` |
| `/`    | `SLASH`    |
| `%`    | `PERCENT`  |
| `==`   | `EQEQ`     |
| `!=`   | `BANGEQ`   |
| `<`    | `LT`       |
| `>`    | `GT`       |
| `<=`   | `LTEQ`     |
| `>=`   | `GTEQ`     |
| `!` | `BANG` |

---

### Keywords

| Lexeme  | TokenType |
| ------- | --------- |
| `let`   | `LET`     |
| `true`  | `TRUE`    |
| `false` | `FALSE`   |
| `null`  | `NULL`    |
| `if`    | `IF`      |
| `else`  | `ELSE`    |

Keywords must be recognized ahead of identifiers.

---

### Identifiers

```ebnf
IDENT = letter { letter | digit | "_" } ;
```

* Must begin with a Unicode letter
* Common examples: `map`, `tile`, `load`

---

### Literals

#### Numbers

```ebnf
NUMBER = digit { digit } [ "." digit { digit } ] ;
```

* Parsed as int64 or float64
* No hex, octal, or exponent support (yet)
* Examples: `42`, `3.14`, `-0.001`

#### Strings

| Delimiter | TokenType |
| --------- | --------- |
| `'`       | `STRING`  |
| `"`       | `STRING`  |

* On valid Go escape sequences allowed: `\n`, `\"`, `\\`, etc.
* Unterminated strings produce an `ILLEGAL` token
* Strings must include the opening and closing delimiter; this is needed by the parser

#### Template Strings

| Delimiter   | TokenType  |
| ----------- | ---------- |
| `` `...` `` | `TEMPLATE` |

* Contain raw text + `${expression}` interpolations
* Must be passed through as a raw lexeme; split into parts during AST construction
* Nested backticks not allowed
* Unterminated template produces `ILLEGAL`

---

### Whitespace & Comments

* Whitespace separates tokens and is otherwise ignored
* Single-line comments begin with `//` and end at newline
* Multiline comments (`/* ... */`) are not supported

---

### Special Tokens

| Purpose       | TokenType |
| ------------- | --------- |
| End of file   | `EOF`     |
| Invalid input | `ILLEGAL` |

---

## 🧪 Error Handling

* Emit `ILLEGAL` for:
    * Unrecognized characters
    * Unterminated strings or template literals
    * Invalid escape sequences
* Always include `Pos` in error tokens
* Lexer should not panic—errors are recoverable in the parser/VM

---

## 🧹 Suggested Lexer Behavior

* Collapse `==`, `!=`, `<=`, `>=` as two-character tokens
* Scan the longest matching operator/punctuation first
* Automatically insert `SEMICOLON` if not found? ❌ No. Require semicolons for clarity

---

## ✅ Example Input → Tokens

```js
let map = load("world.wxx");
if (map != null) {
  print(`Map has ${map.width} tiles wide.`);
}
```

Produces:

```
LET         "let"
IDENT       "map"
EQUAL       "="
IDENT       "load"
LPAREN      "("
STRING      "\"world.wxx\""
RPAREN      ")"
SEMICOLON   ";"
IF          "if"
LPAREN      "("
IDENT       "map"
BANGEQ      "!="
NULL        "null"
RPAREN      ")"
LBRACE      "{"
IDENT       "print"
LPAREN      "("
TEMPLATE    "`Map has ${map.width} tiles wide.`"
RPAREN      ")"
SEMICOLON   ";"
RBRACE      "}"
EOF
```

---

## 🧰 Implementation Tip

Write a debug mode that dumps token stream with positions:

```go
fmt.Printf("%-10s %-20q @ %d:%d\n", tok.Type, tok.Lexeme, tok.Pos.Line, tok.Pos.Column)
```

---

## 📂 Test Cases to Cover

* All punctuation and operators
* All keywords vs identifiers
* Valid and invalid numbers
* All string literal variants (including escapes)
* Template strings with and without interpolations
* Illegal tokens (e.g., `@`, unterminated strings)
* Comments and whitespace

---

## 🛑 Out of Scope for V1

* No `+=`, `--`, `++`, or compound operators
* No hex, octal, or exponent numbers
* No multi-line string or template escapes
* No multiline (`/* ... */`) comments

---

## 🧭 Notes

This lexer will be used by a recursive descent parser, so the token stream should:

* Be greedy (longest match)
* Track all source positions
* Never panic

---
