package parser

import (
	"github.com/hamdan-khan/interpreter/errorHandler"
	"github.com/hamdan-khan/interpreter/syntax"
	"github.com/hamdan-khan/interpreter/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) Parser {
	return Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() (stmtType []syntax.Stmt, e error) {
	statements := []syntax.Stmt{}

	for !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return stmtType, err
		}
		statements = append(statements, dec)
	}

	return statements, nil
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

// returns next token, no side-effect
func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.EOF
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// checks next token for the type passed
func (p *Parser) check(tok token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tok
}

// compares the given token types with the current token
// and advances if a match is found
func (p *Parser) match(toks ...token.TokenType) bool {
	for _, t := range toks {
		if p.check(t) {
			p.advance() // advance in case of a matching token type
			return true
		}
	}
	return false
}

// assignment -> IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() (syntax.Expr, error) {
	// how can left side (l-value) of an assignment be an expression?
	// example: someObject(x+y).someField = 10
	// does this mean any expression can be an assignment target?
	// no, only variables can be assignment targets which we later validate
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		operator := p.previous()
		right, err := p.assignment()

		if err != nil {
			return nil, err
		}

		if expr, ok := expr.(*syntax.Variable); ok {
			return &syntax.Assign{
				Name:  expr.Name,
				Value: right,
			}, nil
		}
		return nil, p.error(operator, "Invalid assignment target.")
	}

	return expr, nil
}

// expression -> assignment ;
func (p *Parser) expression() (syntax.Expr, error) {
	return p.assignment()
}

// equality -> comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() (syntax.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.EQUAL_EQUAL, token.NOT_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &syntax.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() (syntax.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &syntax.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (syntax.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &syntax.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (syntax.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &syntax.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// unary -> ( "!" | "-" ) unary | call
func (p *Parser) unary() (syntax.Expr, error) {
	if p.match(token.EXCLAMATION, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &syntax.Unary{
			Right:    right,
			Operator: operator,
		}, nil
	}

	return p.call()
}

// call -> primary ( "(" arguments? ")" )*
func (p *Parser) call() (syntax.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

// arguments -> expression ( "," expression )*
func (p *Parser) finishCall(callee syntax.Expr) (syntax.Expr, error) {
	args := []syntax.Expr{}

	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(args) >= 255 {
				return nil, p.error(p.peek(), "Too many arguments. (limit = 255)")
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)

			// if comma is not found after argument, it means the args list is consumed
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(token.RIGHT_PAREN, "Expected ')' after arguments")
	if err != nil {
		return nil, err
	}

	return &syntax.Call{
		Callee:    callee,
		Arguments: args,
		Paren:     paren,
	}, nil
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER
func (p *Parser) primary() (syntax.Expr, error) {
	if p.match(token.FALSE) {
		return &syntax.Literal{Value: false}, nil
	}
	if p.match(token.TRUE) {
		return &syntax.Literal{Value: true}, nil
	}
	if p.match(token.NIL) {
		return &syntax.Literal{Value: nil}, nil
	}
	if p.match(token.STRING, token.NUMBER) {
		return &syntax.Literal{Value: p.previous().Literal}, nil
	}
	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expected ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &syntax.Grouping{Expression: expr}, nil
	}
	if p.match(token.IDENTIFIER) {
		return &syntax.Variable{Name: p.previous()}, nil
	}
	return nil, p.error(p.peek(), "Expected expression.")
}

func (p *Parser) error(tok token.Token, message string) error {
	if tok.TokenType == token.EOF {
		return errorHandler.ReportError(tok.LineNumber, "at end", message)
	}
	return errorHandler.ReportError(tok.LineNumber, "at '"+tok.Lexeme+"'", message)
}

// checks and advances if the provided type matches with the next token
func (p *Parser) consume(token token.TokenType, message string) (t token.Token, err error) {
	if p.check(token) {
		return p.advance(), nil
	}

	return t, p.error(p.peek(), message)
}

// to get to the end of the statement where the error has occured
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == token.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case token.FUNCTION:
		case token.VAR:
		case token.FOR:
		case token.IF:
		case token.WHILE:
		case token.PRINT:
		case token.RETURN:
			return
		}

		p.advance()
	}
}

// statements stuff

// declaration -> funcDecl | varDecl | statement
func (p *Parser) declaration() (syntax.Stmt, error) {
	if p.match(token.FUNCTION) {
		f, err := p.function("function")
		if err != nil {
			p.synchronize()
			return nil, err
		}
		return f, nil
	}
	if p.match(token.VAR) {
		v, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, err
		}
		return v, nil
	}

	s, sErr := p.statement()
	if sErr != nil {
		p.synchronize()
		return nil, sErr
	}
	return s, nil
}

// function -> IDENTIFIER "(" parameters? ")" block
func (p *Parser) function(kind string) (syntax.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expected "+kind+" name")
	if err != nil {
		return nil, err
	}
	p.consume(token.LEFT_PAREN, "Expected '(' after "+kind+" name")
	parameters := []token.Token{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, p.error(p.peek(), "Too many parameters. (limit = 255)")
			}
			paramName, err := p.consume(token.IDENTIFIER, "Expected parameter name")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, paramName)
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	p.consume(token.RIGHT_PAREN, "Expected ')' after parameters")

	p.consume(token.LEFT_BRACE, "Expected '{' before "+kind+" body")
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	return &syntax.Function{Name: name, Params: parameters, Body: body}, nil
}

// varDecl -> "var" IDENTIFIER ( "=" expression )? ";"
func (p *Parser) varDeclaration() (syntax.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expected variable name")
	if err != nil {
		return nil, err
	}

	var initializer syntax.Expr = nil
	if p.match(token.EQUAL) {
		init, er := p.expression()
		if er != nil {
			return nil, er
		}
		initializer = init
	}
	_, sErr := p.consume(token.SEMICOLON, "Expected ';' after variable declaration")

	if sErr != nil {
		return nil, sErr
	}
	return &syntax.Var{Name: name, Initializer: initializer}, nil
}

// statement -> exprStmt | ifStmt | printStmt | whileStmt | block | returnStmt
func (p *Parser) statement() (syntax.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.RETURN) {
		return p.returnStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.LEFT_BRACE) {
		blockStatements, err := p.blockStatement()
		if err != nil {
			return nil, err
		}
		return &syntax.Block{Statements: blockStatements}, nil
	}
	if p.match(token.IF) {
		return p.ifStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) returnStatement() (syntax.Stmt, error) {
	keyword := p.previous()
	var value syntax.Expr = nil
	// if next token after return keyword is not a semicolon,
	// it means we're returning some value
	if !p.check(token.SEMICOLON) {
		val, err := p.expression()
		if err != nil {
			return nil, err
		}
		value = val
	}
	_, err := p.consume(token.SEMICOLON, "Expected ';' after return")
	if err != nil {
		return nil, err
	}
	return &syntax.Return{Keyword: keyword, Value: value}, nil
}

// syntax desugaring - converting "for" loop into "while" loop
// for (init; condition; increment) body
func (p *Parser) forStatement() (s syntax.Stmt, e error) {
	_, err := p.consume(token.LEFT_PAREN, "Expected '(' after 'for'")
	if err != nil {
		return nil, err
	}

	// initializer can be nil, expression, or a declaration
	var initializer syntax.Stmt = nil
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	// condition can be nil or an expression
	var condition syntax.Expr = nil
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expected ';' after for loop condition")
	if err != nil {
		return nil, err
	}

	// increment can be nil or an expression
	var increment syntax.Expr = nil
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expected ')' after 'for'")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// if present, add the increment expression at the end of the body
	// which would somewhat look like:
	// {
	// 	bodyStatements;
	// 	incrementExpression;
	// }
	if increment != nil {
		body = &syntax.Block{Statements: []syntax.Stmt{body, &syntax.StatementExpression{Expression: increment}}}
	}

	// if condition is absent, make it true i.e. infinite loop
	if condition == nil {
		condition = &syntax.Literal{Value: true}
	}
	body = &syntax.While{Condition: condition, Body: body} // body is now a while loop

	// if initializer is present, make it a block of initializer + body (which is now a while loop)
	if initializer != nil {
		body = &syntax.Block{Statements: []syntax.Stmt{initializer, body}}
	}

	// the "for" loop after desugaring looks like:
	// {
	// 	initializer;
	// 	while (condition) {
	// 		body;
	// 		increment;
	// 	}
	// }
	return body, nil
}

// whileStmt -> "while" "(" expression ")" statement
func (p *Parser) whileStatement() (s syntax.Stmt, e error) {
	_, err := p.consume(token.LEFT_PAREN, "Expected '(' after 'while'")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expected ')' after condition")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &syntax.While{Condition: condition, Body: body}, nil
}

// block -> "{" declaration* "}"
func (p *Parser) blockStatement() (s []syntax.Stmt, e error) {
	statements := []syntax.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.synchronize()
			return nil, err
		}
		statements = append(statements, stmt)
	}

	_, cErr := p.consume(token.RIGHT_BRACE, "Expected '}' after block")
	if cErr != nil {
		return nil, cErr
	}
	return statements, nil
}

func (p *Parser) printStatement() (s syntax.Stmt, e error) {
	expr, err := p.expression()
	if err != nil {
		return s, err
	}

	_, cErr := p.consume(token.SEMICOLON, "Expected ';' after expression")
	if cErr != nil {
		return s, cErr
	}

	return &syntax.Print{Expression: expr}, nil
}

func (p *Parser) expressionStatement() (s syntax.Stmt, e error) {
	expr, err := p.expression()
	if err != nil {
		return s, err
	}

	_, cErr := p.consume(token.SEMICOLON, "Expected ';' after expression!")
	if cErr != nil {
		return s, cErr
	}

	return &syntax.StatementExpression{Expression: expr}, nil
}

// ifStmt -> "if" "(" expression ")" statement ( "else" statement )?
func (p *Parser) ifStatement() (s syntax.Stmt, e error) {
	_, err := p.consume(token.LEFT_PAREN, "Expected '(' after 'if'")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expected ')' after condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	// in case of nested if's and one else branch the innermost if
	// will be associated with the else branch (dangling else problem)
	var elseBranch syntax.Stmt = nil
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &syntax.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

// logic_or -> logic_and ( "or" logic_and )*
func (p *Parser) or() (syntax.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &syntax.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// logic_and -> equality ( "and" equality )*
func (p *Parser) and() (syntax.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &syntax.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}
