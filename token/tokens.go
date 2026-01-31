package token

type TokenType int

const (
	// single char tokens
	LEFT_PAREN TokenType = iota
    RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// possible doule char tokens
	EXCLAMATION
	EQUAL
	NOT_EQUAL // !=
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// literals
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	OR
	TRUE
	FALSE
	FUNCTION
	RETURN
	FOR
	WHILE
	IF
	ELSE
	VAR
	NIL
	PRINT

	// misc
	EOF
)

type Token struct {
	TokenType TokenType
	Lexeme string
	LineNumber int
	Literal any
}


var ReservedKeywords = map[string]TokenType{
	"and": AND,
    "else": ELSE,
    "false": FALSE,
    "for": FOR,
    "fun": FUNCTION,
    "if": IF,
    "nil": NIL,
    "or": OR,
    "print": PRINT,
    "return": RETURN,
    "true": TRUE,
    "var": VAR,
    "while": WHILE,
}