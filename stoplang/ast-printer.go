package stoplang

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr Expr) string {
	result, _ := expr.Accept(p)

	return result.(string)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) (any, error) {
	return p.parenthesize(expr.Operator.lexeme, expr.Left, expr.Right), nil
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) (any, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}

	return fmt.Sprintf("%v", expr.Value), nil
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) (any, error) {
	return p.parenthesize(expr.Operator.lexeme, expr.Right), nil
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		res, _ := expr.Accept(p)
		builder.WriteString(res.(string))
	}

	builder.WriteString(")")

	return builder.String()
}

func (p *AstPrinter) main(args []string) {
	expr := &Binary{
		&Unary{
			Token{tokenType: MINUS, lexeme: "-", literal: nil, line: 1},
			&Literal{Value: 1},
		},
		Token{tokenType: STAR, lexeme: "*", literal: nil, line: 1},
		&Grouping{
			&Literal{Value: 69.67},
		},
	}

	fmt.Println(expr)
}
