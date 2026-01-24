package main

import "go/types"

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

	// possible doule char tokens
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

	// misc
	EOF
)

type Token struct {
	tokenType TokenType
	lexeme string
	lineNumber int
	literal types.Object
}
