package stoplang

import (
	"fmt"
)

type Interpreter struct{}

func (i *Interpreter) Evaluate(expr Expr) (any, error) {
	switch e := expr.(type) {
	case *Literal:
		return e.Value, nil

	case *Grouping:
		return i.Evaluate(e.Expression)

	case *Unary:
		right, err := i.Evaluate(e.Right)
		if err != nil {
			return nil, err
		}
		switch e.Operator.tokenType {
		case MINUS:
			value, ok := right.(float64)
			if !ok {
				return nil, fmt.Errorf("operand must be a number")
			}
			return -value, nil

		case BANG:
			return !i.isTruthy(right), nil
		}

		return nil, nil

		return right, nil

	case *Binary:
		left, err := i.Evaluate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := i.Evaluate(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.tokenType {
		case MINUS:
			leftNumber, rightNumber, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftNumber - rightNumber, nil

		case SLASH:
			leftNumber, rightNumber, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftNumber / rightNumber, nil

		case STAR:
			leftNumber, rightNumber, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftNumber * rightNumber, nil
		}

		return i.executeBinary(left, e.Operator, right)

	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (i *Interpreter) executeBinary(left any, op Token, right any) (any, error) {

	return nil, nil
}

func (i *Interpreter) isTruthy(obj any) bool {
	if obj == nil {
		return false
	}

	if value, ok := obj.(bool); ok {
		return value
	}

	return true
}

func (i *Interpreter) checkNumberOperands(operator Token, left, right any) (float64, float64, error) {
	leftNumber, okLeft := left.(float64)
	rightNumber, okRight := right.(float64)

	if !okLeft || !okRight {
		return 0, 0, fmt.Errorf("line %d: operands must be number for operator %s", operator.line, operator.lexeme)
	}

	return leftNumber, rightNumber, nil
}
