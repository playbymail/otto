// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package parser

import (
	"github.com/playbymail/otto/wjs/ast"
	"github.com/playbymail/otto/wjs/lexer"
)

type Parser struct {
	tokens []lexer.Token
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Start: p.tokens[0].Pos,
	}
	// TODO: implement parser
	//for p.peek().Type != TokenEOF {
	//	stmt := p.parseStatement()
	//	program.Stmts = append(program.Stmts, stmt)
	//}
	return program
}
