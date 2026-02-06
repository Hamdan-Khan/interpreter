package main

import (
	"strconv"

	"github.com/hamdan-khan/interpreter/errorHandler"
	"github.com/hamdan-khan/interpreter/token"
)

type Scanner struct {
	tokens []token.Token
	lineNumber int	// tracks the line number being scanned
	start int 		// tracks the start of the lexeme
	current int 	// tracks the current char
	source string 	// contents of source file
} 

func NewScanner (source string) Scanner {
	return Scanner{
		lineNumber: 1,
		start: 0,
		current: 0,
		source: source,
	}
}

func (s *Scanner) Scan() {
	s.lineNumber = 1

	// scans every char and tokenize lexemes
	for !s.isAtEnd() {
		// after every token addition, the new bound is where the last token ended
		// Current state is updated in scanToken method
		s.start = s.current 
		s.scanToken()
	}

	// to denote that we've reached the end of file
	eofToken := token.Token{TokenType: token.EOF,
					  Lexeme: "",
					  LineNumber: s.lineNumber,
					  Literal: nil}
	s.tokens = append(s.tokens, eofToken)
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
		// single char tokens
		case '(': s.addToken(token.LEFT_PAREN, nil)
		case ')': s.addToken(token.RIGHT_PAREN, nil)
		case '{': s.addToken(token.LEFT_BRACE, nil)
		case '}': s.addToken(token.RIGHT_BRACE, nil)
		case ',': s.addToken(token.COMMA, nil)
		case '.': s.addToken(token.DOT, nil)
		case '-': s.addToken(token.MINUS, nil)
		case '+': s.addToken(token.PLUS, nil)
		case ';': s.addToken(token.SEMICOLON, nil)
		case '*': s.addToken(token.STAR, nil)

		// double char tokens
		case '=':
			if s.match('='){
				s.addToken(token.EQUAL_EQUAL, nil)
			} else{
				s.addToken(token.EQUAL, nil)
			}
		case '!':
			if s.match('='){
				s.addToken(token.NOT_EQUAL, nil)
			} else{
				s.addToken(token.EXCLAMATION, nil)
			}
		case '>':
			if s.match('='){
				s.addToken(token.GREATER_EQUAL, nil)
			} else{
				s.addToken(token.GREATER, nil)
			}
		case '<':
			if s.match('='){
				s.addToken(token.LESS_EQUAL, nil)
			} else{
				s.addToken(token.LESS, nil)
			}

		case '/':
			// trailing slash can be a single-line comment,
			// multi-line comment or a division operator
			if s.match('/'){
				// advance through the comments without tokenizing anything
				// by looking for the new line or the EOF
				for (s.next() != '\n' && !s.isAtEnd()) {
					s.advance()
				}
			} else if s.match('*') {
				// handles block comments /* */, advance till */ is encountered
				for !s.isAtEnd() && !(s.next() == '*' && s.nextNext() == '/') {
					if s.next() == '\n'{
						s.lineNumber ++
					}
					s.advance()
				}

				if s.isAtEnd() {
					errorHandler.ReportError(s.lineNumber, "", "Unterminated block comment.")
				}

				// consume the closing * and / 
				s.advance()
				s.advance()
			} else {
				s.addToken(token.SLASH, nil)
			}

		// ignore space/tabs
		case ' ':
		case '\t':
		case '\n':
			s.lineNumber++

		default:
			if s.isDigit(char) {
				s.handleNumber()
			} else if s.isAlpha(char) { 
				// an alphanum identifier shouldn't start with a digit
				s.handleIdentifier()
			} else {
				errorHandler.ReportError(s.lineNumber, "", "Unexpected character.")
			}
	}
}

func (s *Scanner) isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' ||
		   c >= 'A' && c <= 'Z' ||
		   c == '_'
}

func (s *Scanner) isAlphaNum(c rune) bool {
	return s.isDigit(c) || s.isAlpha(c)
}

func (s *Scanner) handleIdentifier() {
	for s.isAlphaNum(s.next()) {
		s.advance()
	}
	// the scanned alphanum can either be an identifier defined by user or
	// a reserved keyword like and, or, return, etc.
	text := s.source[s.start:s.current]
	tokType, ok := token.ReservedKeywords[text]
	if !ok {
		tokType = token.IDENTIFIER
	}
	s.addToken(tokType, nil)
}

func (s *Scanner) isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

// scans and tokenizes numbers 
func (s *Scanner) handleNumber() {
	for s.isDigit(s.next()){
		s.advance()
	}

	// if the char after the decimal point is a digit,
	// its the fractional part. In some cases, it can be a method too
	if s.next() == '.' && s.isDigit(s.nextNext()) {
		// pointer moved to the decimal point
		s.advance()

		// advance through the fractional part
		for s.isDigit(s.next()) { 
			s.advance()
		}
	}

	numLiteral, err := strconv.ParseFloat(s.source[s.start: s.current], 64)
	if err != nil {
		errorHandler.ReportError(s.lineNumber, "", "Unexpected number encountered")
	}
	s.addToken(token.NUMBER, numLiteral)
}

// returns the char at the current pointer and increments the curr pointer
func (s *Scanner) advance() rune{
	curr := rune(s.source[s.current])
	s.current++
	return curr
}

// checks the next char against the given value and increments the curr pointer
// only if it matches with the expected value
func (s *Scanner) match(char rune) bool{
	if s.isAtEnd() { return false }
	// why check with current? 
	// because current is already incremented using advance in scanToken
	if rune(s.source[s.current]) != char { return false}

	s.current++
	return true
}

// returns the next char
func (s *Scanner) next() rune{
	if s.isAtEnd() { return '\000'} // null terminator
	return rune(s.source[s.current])
}

// returns the char after the next char
func (s *Scanner) nextNext() rune{
	if s.current + 1 >= len(s.source) { return '\000'} // null terminator
	return rune(s.source[s.current+1])
}

func (s *Scanner) isAtEnd() bool{
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(tokenType token.TokenType, literal any){
	lexeme :=  s.source[s.start: s.current]
	token := token.Token{TokenType: tokenType,
				   Literal: literal,
				   Lexeme: lexeme,
				   LineNumber: s.lineNumber}
	s.tokens = append(s.tokens, token)
}
