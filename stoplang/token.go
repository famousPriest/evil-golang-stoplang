package stoplang

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t *Token) ToString() string {
	return fmt.Sprintf("%v %s %v", t.tokenType, t.lexeme, t.literal)
}
