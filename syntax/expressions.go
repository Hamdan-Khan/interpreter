package syntax

import "github.com/hamdan-khan/interpreter/token"

type Visitor interface {
	VisitBinaryExpr(expr *Binary) (any, error)
	VisitGroupingExpr(expr *Grouping) (any, error)
	VisitLiteralExpr(expr *Literal) (any, error)
	VisitUnaryExpr(expr *Unary) (any, error)
	VisitVariableExpr(expr *Variable) (any, error)
	VisitAssignExpr(expr *Assign) (any, error)
	VisitLogicalExpr(expr *Logical) (any, error)
	VisitCallExpr(expr *Call) (any, error)
}

type Expr interface {
	Accept(visitor Visitor) (any, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
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
	Right    Expr
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

// why assignment as expression?
// design choice, assignment can be statement too (like in python), in our case
// it can be nested inside a larger expression, or cases like a = b = 10
type Assign struct {
	Name  token.Token
	Value Expr
}

func (e *Assign) Accept(visitor Visitor) (any, error) {
	return visitor.VisitAssignExpr(e)
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *Logical) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLogicalExpr(e)
}

type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (e *Call) Accept(visitor Visitor) (any, error) {
	return visitor.VisitCallExpr(e)
}
