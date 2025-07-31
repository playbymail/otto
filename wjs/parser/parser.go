// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package parser

import (
	"strconv"

	"github.com/playbymail/otto/wjs/ast"
	"github.com/playbymail/otto/wjs/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) ParseProgram() *ast.Program {
	if len(p.tokens) == 0 {
		return &ast.Program{Stmts: []ast.Stmt{}}
	}

	program := &ast.Program{
		Start: p.tokens[0].Pos,
		Stmts: []ast.Stmt{},
	}

	for p.peek().Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
	}

	return program
}

// Helper methods for token navigation
func (p *Parser) peek() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() lexer.Token {
	token := p.peek()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return token
}

func (p *Parser) expect(tokenType lexer.TokenType) bool {
	if p.peek().Type == tokenType {
		p.advance()
		return true
	}
	return false
}

// Statement parsing
func (p *Parser) parseStatement() ast.Stmt {
	switch p.peek().Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.IDENT:
		// Could be assignment or expression statement
		return p.parseIdentStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() ast.Stmt {
	start := p.advance().Pos // consume 'let'

	if p.peek().Type != lexer.IDENT {
		return nil // error: expected identifier
	}

	name := &ast.Ident{
		Start: p.peek().Pos,
		Name:  p.advance().Lexeme,
	}

	if !p.expect(lexer.EQUAL) {
		return nil // error: expected '='
	}

	value := p.parseExpression()
	if value == nil {
		return nil
	}

	p.expect(lexer.SEMICOLON) // optional semicolon

	return &ast.LetStmt{
		Start: start,
		Name:  name,
		Value: value,
	}
}

func (p *Parser) parseIdentStatement() ast.Stmt {
	// Lookahead to distinguish assignment from expression
	if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == lexer.EQUAL {
		return p.parseAssignmentStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseAssignmentStatement() ast.Stmt {
	target := p.parseExpression()
	if target == nil {
		return nil
	}

	start := target.Pos()

	if !p.expect(lexer.EQUAL) {
		return nil // error: expected '='
	}

	value := p.parseExpression()
	if value == nil {
		return nil
	}

	p.expect(lexer.SEMICOLON) // optional semicolon

	return &ast.AssignStmt{
		Start:  start,
		Target: target,
		Value:  value,
	}
}

func (p *Parser) parseExpressionStatement() ast.Stmt {
	expr := p.parseExpression()
	if expr == nil {
		return nil
	}

	p.expect(lexer.SEMICOLON) // optional semicolon

	return &ast.ExprStmt{
		Start: expr.Pos(),
		Value: expr,
	}
}

// Expression parsing (operator precedence)
const (
	_ int = iota
	LOWEST
	EQUALS      // ==, !=
	LESSGREATER // > or <
	SUM         // +, -
	PRODUCT     // *, /
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index], obj.field
)

var precedences = map[lexer.TokenType]int{
	lexer.EQEQ:     EQUALS,
	lexer.BANGEQ:   EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.LTEQ:     LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.GTEQ:     LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.ASTERISK: PRODUCT,
	lexer.SLASH:    PRODUCT,
	lexer.PERCENT:  PRODUCT,
	lexer.LPAREN:   CALL,
	lexer.DOT:      INDEX,
	lexer.LBRACK:   INDEX,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peek().Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseExpression() ast.Expr {
	return p.parseExpressionWithPrecedence(LOWEST)
}

func (p *Parser) parseExpressionWithPrecedence(precedence int) ast.Expr {
	left := p.parsePrimaryExpression()
	if left == nil {
		return nil
	}

	for p.peek().Type != lexer.SEMICOLON && p.peek().Type != lexer.EOF && precedence < p.peekPrecedence() {
		switch p.peek().Type {
		case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH, lexer.PERCENT,
			lexer.EQEQ, lexer.BANGEQ, lexer.LT, lexer.LTEQ, lexer.GT, lexer.GTEQ:
			left = p.parseBinaryExpression(left)
		case lexer.LPAREN:
			left = p.parseCallExpression(left)
		case lexer.DOT:
			left = p.parseMemberExpression(left)
		case lexer.LBRACK:
			left = p.parseIndexExpression(left)
		default:
			return left
		}
	}

	return left
}

func (p *Parser) parsePrimaryExpression() ast.Expr {
	switch p.peek().Type {
	case lexer.IDENT:
		return p.parseIdentifier()
	case lexer.NUMBER:
		return p.parseNumberLiteral()
	case lexer.STRING:
		return p.parseStringLiteral()
	case lexer.TEMPLATE:
		return p.parseTemplateLiteral()
	case lexer.TRUE, lexer.FALSE, lexer.NULL:
		return p.parseBooleanOrNullLiteral()
	case lexer.MINUS, lexer.BANG:
		return p.parseUnaryExpression()
	case lexer.LPAREN:
		return p.parseGroupedExpression()
	default:
		return nil // error: unexpected token
	}
}

func (p *Parser) parseIdentifier() ast.Expr {
	token := p.advance()
	return &ast.Ident{
		Start: token.Pos,
		Name:  token.Lexeme,
	}
}

func (p *Parser) parseNumberLiteral() ast.Expr {
	token := p.advance()

	// Try to parse as integer first
	if intValue, err := strconv.Atoi(token.Lexeme); err == nil {
		val := int64(intValue)
		return &ast.NumberLit{
			Start:    token.Pos,
			IntVal:   &val,
			FloatVal: nil,
		}
	}

	// If not an integer, try parsing as float
	if floatValue, err := strconv.ParseFloat(token.Lexeme, 64); err == nil {
		return &ast.NumberLit{
			Start:    token.Pos,
			IntVal:   nil,
			FloatVal: &floatValue,
		}
	}

	return nil // error: invalid number
}

func (p *Parser) parseStringLiteral() ast.Expr {
	token := p.advance()
	return &ast.StringLit{
		Start: token.Pos,
		Value: token.Lexeme,
	}
}

func (p *Parser) parseTemplateLiteral() ast.Expr {
	token := p.advance()
	// For now, treat template as simple string
	// TODO: implement proper template parsing with interpolation
	return &ast.TemplateLit{
		Start: token.Pos,
		Parts: []ast.TemplatePart{
			&ast.TextPart{
				Start: token.Pos,
				Value: token.Lexeme,
			},
		},
	}
}

func (p *Parser) parseBooleanOrNullLiteral() ast.Expr {
	token := p.advance()
	switch token.Type {
	case lexer.TRUE:
		return &ast.BooleanLit{
			Start: token.Pos,
			Value: true,
		}
	case lexer.FALSE:
		return &ast.BooleanLit{
			Start: token.Pos,
			Value: false,
		}
	case lexer.NULL:
		return &ast.NullLit{
			Start: token.Pos,
		}
	default:
		return nil
	}
}

func (p *Parser) parseUnaryExpression() ast.Expr {
	token := p.advance()
	operand := p.parseExpressionWithPrecedence(PREFIX)
	if operand == nil {
		return nil
	}
	return &ast.UnaryExpr{
		Start:    token.Pos,
		Operator: token.Lexeme,
		Operand:  operand,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expr {
	p.advance() // consume '('
	expr := p.parseExpression()
	if !p.expect(lexer.RPAREN) {
		return nil // error: expected ')'
	}
	return expr
}

func (p *Parser) parseBinaryExpression(left ast.Expr) ast.Expr {
	token := p.advance()
	precedence := precedences[token.Type]
	right := p.parseExpressionWithPrecedence(precedence)
	if right == nil {
		return nil
	}
	return &ast.BinaryExpr{
		Start:    left.Pos(),
		Left:     left,
		Operator: token.Lexeme,
		Right:    right,
	}
}

func (p *Parser) parseCallExpression(callee ast.Expr) ast.Expr {
	start := callee.Pos()
	p.advance() // consume '('

	args := []ast.Expr{}
	if p.peek().Type != lexer.RPAREN {
		args = append(args, p.parseExpression())
		for p.expect(lexer.COMMA) {
			args = append(args, p.parseExpression())
		}
	}

	if !p.expect(lexer.RPAREN) {
		return nil // error: expected ')'
	}

	return &ast.CallExpr{
		Start:  start,
		Callee: callee,
		Args:   args,
	}
}

func (p *Parser) parseMemberExpression(object ast.Expr) ast.Expr {
	start := object.Pos()
	p.advance() // consume '.'

	if p.peek().Type != lexer.IDENT {
		return nil // error: expected identifier
	}

	field := &ast.Ident{
		Start: p.peek().Pos,
		Name:  p.advance().Lexeme,
	}

	return &ast.MemberExpr{
		Start:  start,
		Object: object,
		Field:  field,
	}
}

func (p *Parser) parseIndexExpression(target ast.Expr) ast.Expr {
	start := target.Pos()
	p.advance() // consume '['

	index := p.parseExpression()
	if index == nil {
		return nil
	}

	if !p.expect(lexer.RBRACK) {
		return nil // error: expected ']'
	}

	return &ast.IndexExpr{
		Start:  start,
		Target: target,
		Index:  index,
	}
}
