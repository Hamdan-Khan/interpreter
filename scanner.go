package main

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	// single char tokens
	LEFT_PAREN TokenType = iota
    RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	EXCLAMATION
	EQUAL
	NOT_EQUAL // !=
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// literals
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	OR
	TRUE
	FALSE
	FUNCTION
	RETURN
	FOR
	WHILE
	IF
	ELSE
	VAR
	NIL
	PRINT
)

type Token struct {
	tokenType TokenType
	lexeme string
	lineNumber int
	literal any
}


func Scan(input string) {
	tokens := strings.Split(input," ")
	for i:= range tokens {
		fmt.Printf("%v, ", tokens[i])
	}
}