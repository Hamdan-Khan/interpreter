package interpreter

import (
	"fmt"

	"github.com/hamdan-khan/interpreter/token"
)

type RuntimeError struct {
	Token   token.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%v at line %d. %v", e.Token.Lexeme, e.Token.LineNumber, e.Message)
}

func NewRuntimeError(t token.Token, msg string) error {
	return &RuntimeError{Token: t, Message: msg}
}

type Return struct {
	Value any
}

func (e *Return) Error() string {
	return ""
}

func NewReturn(value any) error {
	return &Return{Value: value}
}
