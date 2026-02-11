package syntax

type StatementVisitor interface {
    VisitExpressionStmt(expr *StatementExpression) (any, error)
    VisitPrintStmt(expr *Print) (any, error)
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

