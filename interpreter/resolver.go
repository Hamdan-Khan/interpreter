package interpreter

import "github.com/hamdan-khan/interpreter/syntax"

type Resolver struct {
	interpreter *Interpreter
	scopes      []map[string]bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter}
}

func (r *Resolver) VisitBlockStmt(stmt *syntax.Block) error {
	r.beginScope()
	defer r.endScope()

	if err := r.ResolveStmts(stmt.Statements); err != nil {
		return err
	}
	return nil
}

func (r *Resolver) ResolveStmts(stmts []syntax.Stmt) error {
	for _, stmt := range stmts {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) ResolveExpr(expr syntax.Expr) error {
	return expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	// pops the last element
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) VisitVarStmt(stmt *syntax.Var) error {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		if err := r.ResolveExpr(stmt.Initializer); err != nil {
			return err
		}
	}
	r.define(stmt.Name)
	return nil
}
