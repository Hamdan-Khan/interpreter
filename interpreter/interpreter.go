package interpreter

import (
	"fmt"
	"strings"

	"github.com/hamdan-khan/interpreter/errorHandler"
	"github.com/hamdan-khan/interpreter/syntax"
	"github.com/hamdan-khan/interpreter/token"
)

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(stmts []syntax.Stmt) error {
	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// recursively evaluates given expression to produce a literal
// uses visitor pattern to implement functions for each expressions (todo: clarify)
func (i *Interpreter) evaluate(expr syntax.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt syntax.Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) VisitBlockStmt(stmt *syntax.Block) (any, error) {
	err := i.executeBlock(stmt.Statements, NewEnvironmentWithParent(i.environment))
	return nil, err
}

// executes a block of statements in a new environment
func (i *Interpreter) executeBlock(statements []syntax.Stmt, environment *Environment) error {
	// store the parent environment temporarily
	var prev *Environment = i.environment
	// set the interpreter's environment to the new one for block execution
	i.environment = environment
	for _, stmt := range statements {
		_, err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	// restore the parent environment
	i.environment = prev
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *syntax.Var) (any, error) {
	// if initializer ( "=" expression ) is absent, value is nil
	var val any = nil
	if stmt.Initializer != nil {
		v, err := i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
		val = v
	}
	i.environment.Define(stmt.Name.Lexeme, val)
	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(expr *syntax.Variable) (any, error) {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr *syntax.Assign) (any, error) {
	val, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	err = i.environment.Assign(expr.Name, val)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (i *Interpreter) VisitIfStmt(stmt *syntax.If) (any, error) {
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(condition) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil, nil
}

func (i *Interpreter) VisitPrintStmt(stmt *syntax.Print) (any, error) {
	val, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", i.stringify(val))
	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *syntax.StatementExpression) (any, error) {
	_, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
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
	boolVal, isBool := val.(bool) // type assertion (can use comma ok to validate if its an integer)
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
	switch expr.Operator.TokenType {
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

	switch expr.Operator.TokenType {
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
		return i.isEqual(left, right), nil
	case token.NOT_EQUAL:
		return !i.isEqual(left, right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *syntax.Logical) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	// short circuit evaluation (jump ahead):
	// for "or", if the left operand is truthy, we don't evaluate the right operand
	// for "and", if the left operand is falsy, we don't evaluate the right operand
	if expr.Operator.TokenType == token.OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
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

func (i *Interpreter) isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
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
func (i *Interpreter) checkNumberOperand(operator token.Token, operand any) (float64, error) {
	val, ok := operand.(float64)
	if !ok {
		return 0, errorHandler.NewRuntimeError(operator, "Operator must be a number")
	}
	return val, nil
}

// for binary mathematical evaluation
//
// this raises an evaluation error when operands with wrong types are encountered
// expected types are number
func (i *Interpreter) checkNumberOperands(operator token.Token, left any, right any) (leftVal float64, rightVal float64, err error) {
	leftVal, lOk := left.(float64)
	rightVal, rOk := right.(float64)
	if !lOk || !rOk {
		return 0, 0, errorHandler.NewRuntimeError(operator, "Operator must be a number")
	}
	return leftVal, rightVal, nil
}
