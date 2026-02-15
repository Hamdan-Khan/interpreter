package interpreter

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) (any, error)
}

type NativeCallable struct {
	fn    func([]any) (any, error)
	arity int
}

func (nc *NativeCallable) Arity() int {
	return nc.arity
}

func (nc *NativeCallable) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return nc.fn(arguments)
}
