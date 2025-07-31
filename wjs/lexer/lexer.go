// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package lexer implements a lexical scanner for WJS
package lexer

import "github.com/playbymail/otto/wjs/domain"

type Lexer struct {
	filename string // set only when running from a script
	input    string
	line     int
}

func New(filename, input string) *Lexer {
	return &Lexer{filename: filename, input: input, line: 1}
}

func (l *Lexer) NextToken() Token {
	// TODO: implement
	return Token{Pos: domain.Pos{Script: l.filename, Line: l.line}, Type: EOF}
}
