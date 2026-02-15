package interpreter

import (
	"github.com/hamdan-khan/interpreter/errorHandler"
	"github.com/hamdan-khan/interpreter/token"
)

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(token token.Token) (any, error) {
	if value, ok := e.values[token.Lexeme]; ok {
		return value, nil
	}
	return nil, errorHandler.ReportError(token.LineNumber, "undefined variable", "Variable "+token.Lexeme+" is not defined.")
}

func (e *Environment) Assign(token token.Token, value any) error {
	// assignment cannot create a new variable, the var we're assinging
	// to must be defined first
	if _, ok := e.values[token.Lexeme]; ok {
		e.values[token.Lexeme] = value
		return nil
	}
	//todo: runtime error
	return errorHandler.ReportError(token.LineNumber, "undefined variable", "Variable "+token.Lexeme+" is not defined.")
}
