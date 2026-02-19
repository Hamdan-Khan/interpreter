package interpreter

import (
	"github.com/hamdan-khan/interpreter/errorHandler"
	"github.com/hamdan-khan/interpreter/syntax"
)

type Function struct {
	Declaration *syntax.Function
}

func (f *Function) Call(interpreter *Interpreter, arguments []any) (any, error) {
	// environment local to the called function with all the function's arguments
	env := NewEnvironmentWithParent(interpreter.environment)

	for i, param := range f.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	err := interpreter.executeBlock(f.Declaration.Body, env)
	if err != nil {
		// return disguised as error is used to unwind the stack of statements
		// similar to Java's exception throwing
		if ret, ok := err.(*errorHandler.Return); ok {
			return ret.Value, nil
		}
		// if it's not a return error, it's a runtime error
		return nil, err
	}
	return nil, nil
}

func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

func (f *Function) String() string {
	return "<fn " + f.Declaration.Name.Lexeme + ">"
}
