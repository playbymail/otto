// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := "print(5);"
	
	expected := []struct {
		expectedType   TokenType
		expectedLexeme string
	}{
		{IDENT, "print"},
		{LPAREN, "("},
		{NUMBER, "5"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	l := New("test", input)

	for i, tt := range expected {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
		}
	}
}

func TestAllTokens(t *testing.T) {
	input := "print(5);"
	
	expected := []struct {
		expectedType   TokenType
		expectedLexeme string
	}{
		{IDENT, "print"},
		{LPAREN, "("},
		{NUMBER, "5"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	l := New("test", input)
	tokens := l.AllTokens()

	if len(tokens) != len(expected) {
		t.Fatalf("wrong number of tokens. expected=%d, got=%d", len(expected), len(tokens))
	}

	for i, tt := range expected {
		tok := tokens[i]

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
		}
	}

	// Ensure the last token is EOF
	lastToken := tokens[len(tokens)-1]
	if lastToken.Type != EOF {
		t.Fatalf("last token should be EOF, got=%q", lastToken.Type)
	}
}
