# WSJ Grammar

## 🎯 Grammar (For the Parser Team – EBNF Style)

### ⚠️ **Constraint:** This grammar must support a *recursive descent parser*. So we avoid left recursion and favor LL(1) friendliness.

```ebnf
program     = { statement } ;

statement   = letStmt
            | assignStmt
            | exprStmt ;

letStmt     = LET IDENT EQUAL expression SEMICOLON ;

assignStmt  = target EQUAL expression SEMICOLON ;

target      = IDENT [ DOT IDENT | LBRACK expression RBRACK ] ;

exprStmt    = expression SEMICOLON ;

expression  = equality ;

equality    = comparison [ ( EQEQ | BANGEQ ) comparison ] ;

comparison  = term [ ( LT | GT | LTEQ | GTEQ ) term ] ;

term        = factor { ( PLUS | MINUS ) factor } ;

factor      = unary { ( ASTERISK | SLASH | PERCENT ) unary } ;

unary       = ( MINUS | BANG )? primary ;

primary     = literal
            | IDENT
            | LPAREN expression RPAREN
            | callExpr
            | indexExpr
            | memberExpr ;

callExpr    = IDENT LPAREN [ expressionList ] RPAREN ;

indexExpr   = primary LBRACK expression RBRACK ;

memberExpr  = primary DOT IDENT ;

expressionList = expression { COMMA expression } ;

literal     = NUMBER | STRING | TEMPLATE | TRUE | FALSE | NULL ;
```

---

## 🔍 Notes for the Parser Team

* The parser must **disambiguate primary → callExpr vs memberExpr vs indexExpr** by peeking ahead after `primary`.
* All statements must end in a semicolon (`;`), including expression statements.
* We assume a flat scope for now—no blocks or function declarations.
* The `template` rule is parsed as a single literal by the lexer and split into parts by a post-pass (not recursive descent).
* No need for `return`, `while`, or `function`—this is an expression-oriented scripting language.
