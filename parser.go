package main

import (
	"github.com/hamdan-khan/interpreter/syntax"
	"github.com/hamdan-khan/interpreter/token"
)

type Parser struct {
	tokens []token.Token
	current int
}

func NewParser(tokens []token.Token) Parser {
	return Parser{
		tokens: tokens,
	}
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.EOF
}

func (p *Parser) advance() token.Token {
	if (!p.isAtEnd()) {
		p.current ++
	}
	return p.previous()
}

func (p *Parser) check(tok token.TokenType) bool {
	if (p.isAtEnd()){
		return false
	}
	return p.peek().TokenType == tok
}

// compares the given token types with the current token
// and advances if a match is found
func (p *Parser) match(toks ...token.TokenType) bool {
	for _, t := range toks{
		if p.check(t) {
			p.advance() // advance in case of a matching token type
			return true
		}
	}
	return false
}

func (p *Parser) expression() syntax.Expr {
	return p.equality()
}

// equality -> comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() syntax.Expr {
	expr := p.comparison()

	for p.match(token.EQUAL_EQUAL, token.NOT_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &syntax.Binary{
			Left: expr,
			Operator: operator,
			Right: right,
		}
	}

	return expr
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() syntax.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &syntax.Binary{
			Left: expr,
			Operator: operator,
			Right: right,
		}
	}

	return  expr
}

func (p *Parser) term() syntax.Expr {
	expr := p.factor()

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right := p.factor()
		expr = &syntax.Binary{
			Left: expr,
			Operator: operator,
			Right: right,
		}
	}

	return  expr
}

func (p *Parser) factor() syntax.Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &syntax.Binary{
			Left: expr,
			Operator: operator,
			Right: right,
		}
	}

	return  expr
}

// unary -> ( "!" | "-" ) unary | primary
func (p *Parser) unary() syntax.Expr {
	for p.match(token.EXCLAMATION, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return  &syntax.Unary{
			Right: right,
			Operator: operator,
		}
	}

	return p.primary()
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *Parser) primary() syntax.Expr {
	if p.match(token.FALSE) {
		return &syntax.Literal{Value: false}
	}
	if p.match(token.TRUE) {
		return &syntax.Literal{Value: true}
	}
	if p.match(token.NIL) {
		return &syntax.Literal{Value: nil}
	}
	if p.match(token.STRING, token.NUMBER) {
		return &syntax.Literal{Value: p.previous().Literal}
	}
	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expected ')' after expression.")
		return &syntax.Grouping{Expression: expr}
	}
}

func (p *Parser) consume(token token.TokenType, message string) {
	// todo
}