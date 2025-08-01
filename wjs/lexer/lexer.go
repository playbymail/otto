// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package lexer implements a lexical scanner for WJS
package lexer

import (
	"github.com/playbymail/otto/wjs/domain"
)

type Lexer struct {
	script   string // set only when running from a script
	input    string
	position int  // current position in input (points to current char)
	readPos  int  // current reading position in input (after current char)
	ch       byte // current char under examination
	line     int  // current line
	column   int  // current column
}

func New(script, input string) *Lexer {
	l := &Lexer{script: script, input: input, line: 1, column: 1}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Pos = domain.Pos{
		Script: l.script,
		Line:   l.line,
		Column: l.column,
		Offset: l.position,
	}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = EQEQ
			tok.Lexeme = string(ch) + string(l.ch)
		} else {
			tok.Type = EQUAL
			tok.Lexeme = string(l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = BANGEQ
			tok.Lexeme = string(ch) + string(l.ch)
		} else {
			tok.Type = BANG
			tok.Lexeme = string(l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = LTEQ
			tok.Lexeme = string(ch) + string(l.ch)
		} else {
			tok.Type = LT
			tok.Lexeme = string(l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = GTEQ
			tok.Lexeme = string(ch) + string(l.ch)
		} else {
			tok.Type = GT
			tok.Lexeme = string(l.ch)
		}
	case '+':
		tok.Type = PLUS
		tok.Lexeme = string(l.ch)
	case '-':
		tok.Type = MINUS
		tok.Lexeme = string(l.ch)
	case '*':
		tok.Type = ASTERISK
		tok.Lexeme = string(l.ch)
	case '/':
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken() // Get next token after comment
		} else {
			tok.Type = SLASH
			tok.Lexeme = string(l.ch)
		}
	case '%':
		tok.Type = PERCENT
		tok.Lexeme = string(l.ch)
	case ',':
		tok.Type = COMMA
		tok.Lexeme = string(l.ch)
	case ';':
		tok.Type = SEMICOLON
		tok.Lexeme = string(l.ch)
	case ':':
		tok.Type = COLON
		tok.Lexeme = string(l.ch)
	case '(':
		tok.Type = LPAREN
		tok.Lexeme = string(l.ch)
	case ')':
		tok.Type = RPAREN
		tok.Lexeme = string(l.ch)
	case '[':
		tok.Type = LBRACK
		tok.Lexeme = string(l.ch)
	case ']':
		tok.Type = RBRACK
		tok.Lexeme = string(l.ch)
	case '{':
		tok.Type = LBRACE
		tok.Lexeme = string(l.ch)
	case '}':
		tok.Type = RBRACE
		tok.Lexeme = string(l.ch)
	case '.':
		tok.Type = DOT
		tok.Lexeme = string(l.ch)
	case '"':
		tok.Type = STRING
		tok.Lexeme = l.readString('"')
	case '\'':
		tok.Type = STRING
		tok.Lexeme = l.readString('\'')
	case '`':
		tok.Type = TEMPLATE
		tok.Lexeme = l.readTemplate()
	case 0:
		tok.Lexeme = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Lexeme = l.readIdentifier()
			tok.Type = LookupIdent(tok.Lexeme)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Lexeme = l.readNumber()
			return tok
		} else {
			tok.Type = ILLEGAL
			tok.Lexeme = string(l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.position = l.readPos
	l.readPos++

	if l.ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipLineComment() {
	// Skip the "//"
	l.readChar()
	l.readChar()
	
	// Skip everything until end of line or end of input
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	// Handle decimal numbers
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString(delimiter byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == delimiter || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar() // Skip escape character and next character
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readTemplate() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '`' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar() // Skip escape character and next character
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// AllTokens returns all tokens in the input as a slice, ending with EOF.
func (l *Lexer) AllTokens() []Token {
	var tokens []Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}

	return tokens
}
