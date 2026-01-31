package main

import "github.com/hamdan-khan/interpreter/token"

type Parser struct {
	tokens []token.Token
	current int
}

func NewParser(tokens []token.Token) Parser {
	return Parser{
		tokens: tokens,
	}
}