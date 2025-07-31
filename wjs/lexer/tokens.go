// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lexer

import (
	"fmt"
	"github.com/playbymail/otto/wjs/domain"
)

// TokenType is the type of lexical tokens.
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Identifiers and literals
	IDENT    // main, foo, tile
	NUMBER   // 42, 3.14
	STRING   // "hello", 'world'
	TEMPLATE // `hello ${x}`

	// Operators
	PLUS     // +
	MINUS    // -
	ASTERISK // *
	SLASH    // /
	PERCENT  // %

	EQEQ   // ==
	BANGEQ // !=
	LT     // <
	LTEQ   // <=
	GT     // >
	GTEQ   // >=

	EQUAL // =

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :

	LPAREN // (
	RPAREN // )
	LBRACK // [
	RBRACK // ]
	LBRACE // {
	RBRACE // }

	DOT // .

	// Keywords
	LET
	TRUE
	FALSE
	NULL
	IF
	ELSE
)

// Keywords maps identifier strings to token types
var keywords = map[string]TokenType{
	"let":   LET,
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
	"if":    IF,
	"else":  ELSE,
}

// LookupIdent returns the token type for a given identifier or keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// Token represents a lexical token with position information.
type Token struct {
	Type   TokenType
	Lexeme string
	Pos    domain.Pos // position in the source file
}

func (t Token) String() string {
	if t.Type == IDENT || t.Type == NUMBER || t.Type == STRING || t.Type == TEMPLATE {
		return fmt.Sprintf("%s(%q)", t.Type.String(), t.Lexeme)
	}
	return t.Type.String()
}

// String returns a string representation of the token type.
func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case TEMPLATE:
		return "TEMPLATE"
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case ASTERISK:
		return "*"
	case SLASH:
		return "/"
	case PERCENT:
		return "%"
	case EQEQ:
		return "=="
	case BANGEQ:
		return "!="
	case LT:
		return "<"
	case LTEQ:
		return "<="
	case GT:
		return ">"
	case GTEQ:
		return ">="
	case EQUAL:
		return "="
	case COMMA:
		return ","
	case SEMICOLON:
		return ";"
	case COLON:
		return ":"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACK:
		return "["
	case RBRACK:
		return "]"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case DOT:
		return "."
	case LET:
		return "let"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case NULL:
		return "null"
	case IF:
		return "if"
	case ELSE:
		return "else"
	default:
		return fmt.Sprintf("TokenType(%d)", int(tt))
	}
}
