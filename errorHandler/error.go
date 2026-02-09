package errorHandler

import (
	"fmt"

	"github.com/hamdan-khan/interpreter/token"
)

func ReportError(lineNumber int, location string, errorMessage string) error {
	fmt.Printf("[line %d] Error %s: %s \n", lineNumber, location, errorMessage)
	return fmt.Errorf("%s", errorMessage)
}

type RuntimeError struct {
	Token   token.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return e.Message
}

func NewRuntimeError(t token.Token, msg string) error {
	return &RuntimeError{Token: t, Message: msg}
}