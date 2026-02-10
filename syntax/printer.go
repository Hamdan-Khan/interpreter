package syntax

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (p *AstPrinter) Print(expr Expr) (string, error) {
	result, err := expr.Accept(p)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) (any, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

// parenthesize wraps expressions in Lisp-style parentheses
// for example: parenthesize("+", left, right) produces "(+ left right)"
func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		result, err := expr.Accept(p)
		if err != nil {
			builder.WriteString("<error>")
			continue
		}
		builder.WriteString(result.(string))
	}

	builder.WriteString(")")
	return builder.String()
}