package syntax

import "github.com/hamdan-khan/interpreter/token"

type Visitor interface {
    VisitBinaryExpr(expr *Binary) (any, error)
    VisitGroupingExpr(expr *Grouping) (any, error)
    VisitLiteralExpr(expr *Literal) (any, error)
    VisitUnaryExpr(expr *Unary) (any, error)
    VisitVariableExpr(expr *Variable) (any, error)
}

type Expr interface {
    Accept(visitor Visitor) (any, error)
}

type Binary struct {
    Left Expr
    Operator token.Token
    Right Expr
}

func (e *Binary) Accept(visitor Visitor) (any, error) {
    return visitor.VisitBinaryExpr(e)
}

type Grouping struct {
    Expression Expr
}

func (e *Grouping) Accept(visitor Visitor) (any, error) {
    return visitor.VisitGroupingExpr(e)
}

type Literal struct {
    Value any
}

func (e *Literal) Accept(visitor Visitor) (any, error) {
    return visitor.VisitLiteralExpr(e)
}

type Unary struct {
    Operator token.Token
    Right Expr
}

func (e *Unary) Accept(visitor Visitor) (any, error) {
    return visitor.VisitUnaryExpr(e)
}

type Variable struct {
    Name token.Token
}

func (e *Variable) Accept(visitor Visitor) (any, error) {
    return visitor.VisitVariableExpr(e)
}