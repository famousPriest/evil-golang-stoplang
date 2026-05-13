package stoplang

import (
	"fmt"
	"slices"
)

type Parser struct {
	tokens  []Token
	current int
	lang    *Lang
}

func (p *Parser) Parse() Expr {
	expr, error := p.expression()
	if error != nil {
		return nil
	}

	return expr
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, error := p.comparison()
	if error != nil {
		return nil, fmt.Errorf("oopsie")
	}
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, error := p.comparison()
		if error != nil {
			return nil, fmt.Errorf("oopsie")
		}

		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, error := p.term()
	if error != nil {
		return nil, fmt.Errorf("oopsie")
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, error := p.term()
		if error != nil {
			return nil, fmt.Errorf("oopsie")
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, error := p.factor()

	if error != nil {
		return nil, fmt.Errorf("oopsie")
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, error := p.factor()
		if error != nil {
			return nil, fmt.Errorf("oopsie")
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, error := p.unary()
	if error != nil {
		return nil, fmt.Errorf("oopsie")
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, error := p.unary()
		if error != nil {
			return nil, fmt.Errorf("oopsie")
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(SLASH, STAR) {
		operator := p.previous()
		right, error := p.unary()
		if error != nil {
			return nil, fmt.Errorf("oopsie")
		}
		return &Binary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Literal{
			Value: false,
		}, nil
	}

	if p.match(TRUE) {
		return &Literal{
			Value: true,
		}, nil
	}

	if p.match(NIL) {
		return &Literal{
			Value: nil,
		}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Literal{
			Value: p.previous().literal,
		}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, error := p.expression()
		if error != nil {
			return nil, fmt.Errorf("oopsie")
		}
		p.consume(RIGHT_PAREN, "expect ')' after expression")

		return &Grouping{
			Expression: expr,
		}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return Token{}, p.error(p.peek(), message)
}

func (p *Parser) error(token Token, message string) error {
	p.lang.ErrorForParser(token, message)

	return fmt.Errorf("parse error")
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS:
		case FUNCTION:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(types ...TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	if p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
