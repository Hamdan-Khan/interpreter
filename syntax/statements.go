package syntax

import "github.com/hamdan-khan/interpreter/token"

type StatementVisitor interface {
	VisitExpressionStmt(expr *StatementExpression) (any, error)
	VisitPrintStmt(expr *Print) (any, error)
	VisitVarStmt(expr *Var) (any, error)
	VisitBlockStmt(expr *Block) (any, error)
	VisitIfStmt(expr *If) (any, error)
}

type Stmt interface {
	Accept(visitor StatementVisitor) (any, error)
}

type StatementExpression struct {
	Expression Expr
}

func (e *StatementExpression) Accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitExpressionStmt(e)
}

type Print struct {
	Expression Expr
}

func (e *Print) Accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitPrintStmt(e)
}

type Var struct {
	Name        token.Token
	Initializer Expr
}

func (e *Var) Accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitVarStmt(e)
}

type Block struct {
	Statements []Stmt
}

func (e *Block) Accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitBlockStmt(e)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e *If) Accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitIfStmt(e)
}
