package stoplang

import (
	"strconv"
)

type Scanner struct {
	lang     *Lang
	source   []rune
	tokens   []Token
	start    int
	current  int
	line     int
	keywords map[string]TokenType
}

// TODO: might want to rename a few keywords
var keywords = map[string]TokenType{
	"and":      AND,
	"class":    CLASS,
	"else":     ELSE,
	"false":    FALSE,
	"for":      FOR,
	"function": FUNCTION,
	"if":       IF,
	"nil":      NIL,
	"or":       OR,
	"print":    PRINT,
	"return":   RETURN,
	"super":    SUPER,
	"self":     SELF,
	"true":     TRUE,
	"var":      VAR,
	"while":    WHILE,
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  []rune(source),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{
		tokenType: EOF,
		lexeme:    "",
		literal:   nil,
		line:      s.line,
	})

	return s.tokens
}

func (s *Scanner) scanToken() {
	r := s.advance()

	switch r {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	case 'o':
		if s.peek() == 'r' {
			s.addToken(OR)
		}
	default:
		if s.isDigit(r) {
			s.number()
		} else if s.isAlpha(r) {
			s.identifier()
		} else {
			s.lang.Error(s.line, "unexpected error")
		}
	}
}

func (s *Scanner) advance() rune {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	text := string(s.source[s.start:s.current])

	s.tokens = append(s.tokens, Token{
		tokenType: tokenType,
		lexeme:    text,
		literal:   literal,
		line:      s.line,
	})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}

	return s.source[s.current]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.lang.Error(s.line, "unterminated string")
		return
	}

	s.advance()

	v := s.source[s.start+1 : s.current-1]

	s.addTokenWithLiteral(STRING, v)
}

func (s *Scanner) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	v, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		s.lang.Error(s.line, "invalid number")
		return
	}

	s.addTokenWithLiteral(NUMBER, v)
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}

	return s.source[s.current+1]

}

func (s *Scanner) isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func (s *Scanner) isAlphaNumeric(r rune) bool {
	return s.isAlpha(r) || s.isDigit(r)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])

	tokenType, isKeyword := keywords[text]

	if !isKeyword {
		tokenType = IDENTIFIER
	}

	s.addToken(tokenType)
}
