package main

import "go/types"

type Scanner struct {
	tokens []Token
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
	eofToken := Token{tokenType: EOF,
					  lexeme: "",
					  lineNumber: s.lineNumber,
					  literal: nil}
	s.tokens = append(s.tokens, eofToken)
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
		// single char tokens
		case '(': s.addToken(LEFT_PAREN, nil)
		case ')': s.addToken(RIGHT_PAREN, nil)
		case '{': s.addToken(LEFT_BRACE, nil)
		case '}': s.addToken(RIGHT_BRACE, nil)
		case ',': s.addToken(COMMA, nil)
		case '.': s.addToken(DOT, nil)
		case '-': s.addToken(MINUS, nil)
		case '+': s.addToken(PLUS, nil)
		case ';': s.addToken(SEMICOLON, nil)
		case '*': s.addToken(STAR, nil)

		// double char tokens
		case '=':
			if s.match('='){
				s.addToken(EQUAL_EQUAL, nil)
			} else{
				s.addToken(EQUAL, nil)
			}
		case '!':
			if s.match('='){
				s.addToken(NOT_EQUAL, nil)
			} else{
				s.addToken(EXCLAMATION, nil)
			}
		case '>':
			if s.match('='){
				s.addToken(GREATER_EQUAL, nil)
			} else{
				s.addToken(GREATER, nil)
			}
		case '<':
			if s.match('='){
				s.addToken(LESS_EQUAL, nil)
			} else{
				s.addToken(LESS, nil)
			}
		
		case '/':
			// trailing slash can be a comment, otherwise it's a division operator
			if s.match('/'){
				// advance through the comments without tokenizing anything
				// by looking for the new line or the EOF
				for (s.next() != '\n' && !s.isAtEnd()) {
					s.advance()
				}
			} else {
				s.addToken(SLASH, nil)
			}
		
		// ignore space/tabs
		case ' ':
		case '\t':
		case '\n':
			s.lineNumber++
		
		default: ReportError(s.lineNumber, "", "Unexpected character.")
	}
}

// returns the char at the current pointer and increments the curr pointer
func (s *Scanner) advance() rune{
	curr := rune(s.source[s.current])
	s.current++
	return curr
}

// checks the next char against the given value
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

func (s *Scanner) isAtEnd() bool{
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(tokenType TokenType, literal types.Object){
	lexeme :=  s.source[s.start: s.current]
	token := Token{tokenType: tokenType,
				   literal: literal,
				   lexeme: lexeme,
				   lineNumber: s.lineNumber}
	s.tokens = append(s.tokens, token)
}
