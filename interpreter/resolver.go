package interpreter

import (
	"github.com/hamdan-khan/interpreter/errorHandler"
	"github.com/hamdan-khan/interpreter/syntax"
	"github.com/hamdan-khan/interpreter/token"
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, currentFunction: NONE}
}

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)

func (r *Resolver) VisitBlockStmt(stmt *syntax.Block) (any, error) {
	r.beginScope()
	defer r.endScope()

	if err := r.ResolveStmts(stmt.Statements); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) ResolveStmts(stmts []syntax.Stmt) error {
	for _, stmt := range stmts {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt syntax.Stmt) error {
	_, err := stmt.Accept(r)
	return err
}

func (r *Resolver) resolveExpr(expr syntax.Expr) error {
	_, err := expr.Accept(r)
	return err
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	// pops the last element
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) VisitVarStmt(stmt *syntax.Var) (any, error) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		if err := r.resolveExpr(stmt.Initializer); err != nil {
			return nil, err
		}
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) declare(name token.Token) {
	// if we are not in a scope, we don't need to declare
	if len(r.scopes) == 0 {
		return
	}
	// if the variable is already declared in the current scope, throw an error
	// can't have two variables with the same name in the same local scope
	if _, ok := r.scopes[len(r.scopes)-1][name.Lexeme]; ok {
		errorHandler.ReportError(name.LineNumber, "", "Already variable with this name in this scope.")
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = false
}

func (r *Resolver) define(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) VisitVariableExpr(expr *syntax.Variable) (any, error) {
	if len(r.scopes) != 0 {
		if defined, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; ok && !defined {
			err := errorHandler.ReportError(expr.Name.LineNumber, "", "Cannot read local variable in its own initializer.")
			return nil, err
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) resolveLocal(expr syntax.Expr, name token.Token) {
	// start at the innermost scope and work outwards
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) VisitAssignExpr(expr *syntax.Assign) (any, error) {
	if err := r.resolveExpr(expr.Value); err != nil {
		return nil, err
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *syntax.Function) (any, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	if err := r.resolveFunction(stmt, FUNCTION); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) resolveFunction(function *syntax.Function, functionType FunctionType) error {
	parentFunction := r.currentFunction
	r.currentFunction = functionType
	defer func() {
		r.currentFunction = parentFunction
	}()
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	if err := r.ResolveStmts(function.Body); err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *syntax.StatementExpression) (any, error) {
	if err := r.resolveExpr(stmt.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt *syntax.If) (any, error) {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}
	if err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		if err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt *syntax.Print) (any, error) {
	if err := r.resolveExpr(stmt.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt *syntax.Return) (any, error) {
	if r.currentFunction == NONE {
		err := errorHandler.ReportError(stmt.Keyword.LineNumber, "", "Cannot return from top-level code.")
		return nil, err
	}
	if stmt.Value != nil {
		if err := r.resolveExpr(stmt.Value); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *syntax.While) (any, error) {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}
	if err := r.resolveStmt(stmt.Body); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *syntax.Binary) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *syntax.Call) (any, error) {
	if err := r.resolveExpr(expr.Callee); err != nil {
		return nil, err
	}

	for _, argument := range expr.Arguments {
		if err := r.resolveExpr(argument); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *syntax.Grouping) (any, error) {
	if err := r.resolveExpr(expr.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr *syntax.Literal) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *syntax.Logical) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *syntax.Unary) (any, error) {
	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}
