package interpreter

import (
	"fmt"
	"strings"

	"github.com/hamdan-khan/interpreter/errorHandler"
	"github.com/hamdan-khan/interpreter/syntax"
	"github.com/hamdan-khan/interpreter/token"
)

type Interpreter struct {}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr syntax.Expr) {
	val, err := i.evaluate(expr)
	if err != nil {

	}
	fmt.Printf("%v\n", i.stringify(val))
}

// recursively evaluates given expression to produce a literal
// uses visitor pattern to implement functions for each expressions (todo: clarify)
func (i *Interpreter) evaluate(expr syntax.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expr *syntax.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *syntax.Grouping) (any, error) {
	return i.evaluate(expr.Expression)
}

// decides which values do we consider to be truthy
// e.g. an empty string, 0, "0", empty array, etc. must have a truthy
// value i.e. true / false during evaluation
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

func (i *Interpreter) VisitUnaryExpr(expr *syntax.Unary) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	// post order traversal (left -> right subtree of AST)
	switch (expr.Operator.TokenType) {
		case token.MINUS:
			val, err := i.checkNumberOperand(expr.Operator, right)
			if err != nil {
				return nil, err
			}
			return -val, nil
		case token.EXCLAMATION:
			return !i.isTruthy(right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *syntax.Binary) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch (expr.Operator.TokenType) {
		case token.MINUS:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftVal - rightVal, nil
		case token.SLASH:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftVal / rightVal, nil
		case token.STAR:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftVal * rightVal, nil
		case token.PLUS:
			return i.executeAdd(left, right), nil
		case token.GREATER:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftVal > rightVal, nil
		case token.GREATER_EQUAL:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftVal >= rightVal, nil
		case token.LESS:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftVal < rightVal, nil
		case token.LESS_EQUAL:
			leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftVal <= rightVal, nil
		case token.EQUAL_EQUAL:
			return i.isEqual(left,right), nil
		case token.NOT_EQUAL:
			return !i.isEqual(left,right), nil
	}

	return nil, nil
}

// executes binary expressions with + operator depending on
// the type i.e. concatenate for strings, add for numbers
func (i *Interpreter) executeAdd(left any, right any) any {
	switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return l + r
			}
		case string:
			if r, ok := right.(string); ok {
				return l + r
			}
	}
	return fmt.Errorf("Invalid operands type for + operation")
}

func (i *Interpreter) isEqual(a any, b any) bool{
	if (a == nil && b == nil){ 
		return true
	}
	if (a == nil) {
		return false
	}

	return a == b
}

func (i *Interpreter) stringify(value any) string {
	if value == nil {
		return "nil"
	}
	switch v := value.(type) {
		case float64:
			text := fmt.Sprintf("%g", v)
			text = strings.TrimSuffix(text, ".0")
			return text
	}

	return fmt.Sprintf("%v", value)
}

// for unary mathematical evaluation
// 
// this raises an evaluation error when operand with wrong type is encountered.
// expected type is number
func (i *Interpreter) checkNumberOperand(operator token.Token, operand any) (float64 ,error){
	val, ok := operand.(float64)
	if (!ok) {
		return 0, errorHandler.NewRuntimeError(operator, "Operator must be a number")
	}
	return val, nil
}


// for binary mathematical evaluation
// 
// this raises an evaluation error when operands with wrong types are encountered
// expected types are number
func (i *Interpreter) checkNumberOperands(operator token.Token, left any, right any) (leftVal float64, rightVal float64, err error){
	leftVal, lOk := left.(float64)
	rightVal, rOk := right.(float64)
	if (!lOk || !rOk) {
		return 0, 0, errorHandler.NewRuntimeError(operator, "Operator must be a number")
	}
	return leftVal, rightVal, nil
}
