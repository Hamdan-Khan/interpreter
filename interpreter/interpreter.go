package interpreter

import (
	"fmt"

	"github.com/hamdan-khan/interpreter/errorHandler"
	"github.com/hamdan-khan/interpreter/syntax"
	"github.com/hamdan-khan/interpreter/token"
)

type Interpreter struct {}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

// recursively evaluates given expression to produce a literal
// uses visitor pattern to implement functions for each expressions (todo: clarify)
func (i *Interpreter) evaluate(expr syntax.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expr *syntax.Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *syntax.Grouping) any {
	return i.evaluate(expr.Expression)
}

// decides which values do we consider to be truthy
// e.g. an empty string, 0, "0", empty array, etc.
func (i *Interpreter) isTruthy(val any) bool {
	if val == nil {
		return false
	}
	boolVal, isBool := val.(bool); // type assertion (can use comma ok to validate if its an integer)
	if isBool {
		return boolVal
	}
	return true
}

func (i *Interpreter) VisitUnaryExpr(expr *syntax.Unary) any {
	right := i.evaluate(expr.Right)
	
	// post order traversal (left -> right subtree of AST)
	switch (expr.Operator.TokenType) {
		case token.MINUS:
			val, err := i.checkNumberOperand(expr.Operator, right)
			if err != nil {
				return err
			}
			return -val
		case token.EXCLAMATION:
			return !i.isTruthy(right)
	}

	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *syntax.Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch (expr.Operator.TokenType) {
		case token.MINUS:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return err
			}
			return leftVal - rightVal
		case token.SLASH:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return err
			}
			return leftVal / rightVal
		case token.STAR:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return err
			}
			return leftVal * rightVal
		case token.PLUS:
			return i.executeAdd(left, right)
		case token.GREATER:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return err
			}
			return leftVal > rightVal
		case token.GREATER_EQUAL:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return err
			}
			return leftVal >= rightVal
		case token.LESS:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return err
			}
			return leftVal < rightVal
		case token.LESS_EQUAL:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return err
			}
			return leftVal <= rightVal
		case token.EQUAL_EQUAL:
			return isEqual(left,right)
		case token.NOT_EQUAL:
			return !isEqual(left,right)
	}

	return nil
}

// executes binary expressions with + operator depending on
// the type i.e. concatenate for strings, add for numbers
func (i *Interpreter) executeAdd(left any, right any) any {
	switch l := left.(type) {
		case int:
			if r, ok := right.(int); ok {
				return l + r
		}
		case string:
			if r, ok := right.(string); ok {
				return l + r
		}
	}
	return fmt.Errorf("Invalid operands type for + operation")
}

func isEqual(a any, b any) bool{
	if (a == nil && b == nil){ 
		return true
	}
	if (a == nil) {
		return false
	}

	return a == b
}

func (i *Interpreter) checkNumberOperand(operator token.Token, operand any) (int ,error){
	val, ok := operand.(int)
	if (!ok) {
		return 0, errorHandler.NewRuntimeError(operator, "Operator must be a number")
	}
	return val, nil
}

func (i *Interpreter) checkNumberOperands(operator token.Token, left any, right any) (leftVal int, rightVal int, err error){
	leftVal, lOk := left.(int)
	rightVal, rOk := right.(int)
	if (!lOk || !rOk) {
		return 0, 0, errorHandler.NewRuntimeError(operator, "Operator must be a number")
	}
	return leftVal, rightVal, nil
}

// todo:
// propagate error upwards to interpreter from visitor functions
// implement error flag in base i/o function