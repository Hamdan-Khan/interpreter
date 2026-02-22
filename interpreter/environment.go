package interpreter

import (
	"github.com/hamdan-khan/interpreter/token"
)

type Environment struct {
	values map[string]any
	// parent-pointer tree (for parent environemnt scope)
	// also called the "cactus stack"
	parent *Environment
}

// for global scope
func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
		parent: nil,
	}
}

// for local scope
func NewEnvironmentWithParent(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]any),
		parent: parent,
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(token token.Token) (any, error) {
	if value, ok := e.values[token.Lexeme]; ok {
		return value, nil
	}
	// local scopes can check parent scopes for variables recursively
	if e.parent != nil {
		return e.parent.Get(token)
	}
	return nil, NewRuntimeError(token, "Undefined variable")
}

func (e *Environment) Assign(token token.Token, value any) error {
	// assignment cannot create a new variable, the var we're assinging
	// to must be defined first
	if _, ok := e.values[token.Lexeme]; ok {
		e.values[token.Lexeme] = value
		return nil
	}
	// local scopes can check parent scopes for variables recursively
	if e.parent != nil {
		return e.parent.Assign(token, value)
	}
	return NewRuntimeError(token, "Undefined variable!")
}

func (e *Environment) GetAt(distance int, name string) (any, error) {
	return e.ancestor(distance).values[name], nil
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	// move up the parent chain "distance" times
	for range distance {
		env = env.parent
	}
	return env
}

func (e *Environment) AssignAt(distance int, name token.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}
