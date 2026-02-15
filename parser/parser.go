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

// assignment -> IDENTIFIER "=" assignment | equality ;
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

// unary -> ( "!" | "-" ) unary | primary
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

	return p.primary()
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

// declaration -> varDecl | statement ;
func (p *Parser) declaration() (syntax.Stmt, error) {
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

// varDecl -> "var" IDENTIFIER ( "=" expression )? ";" ;
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

func (p *Parser) statement() (syntax.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
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

	_, cErr := p.consume(token.SEMICOLON, "Expected ';' after expression")
	if cErr != nil {
		return s, cErr
	}

	return &syntax.StatementExpression{Expression: expr}, nil
}
