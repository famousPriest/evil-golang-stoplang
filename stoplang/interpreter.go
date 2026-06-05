package stoplang

import (
	"fmt"
	"os"
	"strings"
)

type Interpreter struct{}

func (i *Interpreter) Interpret(expression Expr, lang *Lang) {
	value, err := i.Evaluate(expression)
	if err != nil {
		if runtimeErr, ok := err.(*RuntimeError); ok {
			fmt.Fprintf(os.Stderr, "%s\n[line %d]\n", runtimeErr.Message, runtimeErr.Token.line)
			lang.hadRuntimeError = true
		} else {
			fmt.Fprintf(os.Stderr, "Unknown error: %v\n", err)
		}
		return
	}

	fmt.Println(i.stringify(value))
}

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
		case GREATER:
			leftNumber, rightNumber, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftNumber > rightNumber, nil

		case GREATER_EQUAL:
			leftNumber, rightNumber, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftNumber >= rightNumber, nil

		case LESS:
			leftNumber, rightNumber, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftNumber < rightNumber, nil

		case LESS_EQUAL:
			leftNumber, rightNumber, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftNumber <= rightNumber, nil

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

		case PLUS:
			switch l := left.(type) {
			case float64:
				if r, ok := right.(float64); ok {
					return l + r, nil
				}
			case string:
				if r, ok := right.(string); ok {
					return l + r, nil
				}
			}
			return nil, &RuntimeError{
				Token:   e.Operator,
				Message: fmt.Sprintf("line %d: operands must be two numbers or two strings for operator '+'", e.Operator.line),
			}

		case BANG_EQUAL:
			return !i.isEqual(left, right), nil

		case EQUAL_EQUAL:
			return i.isEqual(left, right), nil
		}

		return nil, fmt.Errorf("unknown binary operator: %v", e.Operator)

	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
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

func (i *Interpreter) isEqual(left any, right any) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return left == right
}

func (i *Interpreter) checkNumberOperands(operator Token, left, right any) (float64, float64, error) {
	leftNumber, okLeft := left.(float64)
	rightNumber, okRight := right.(float64)

	if !okLeft || !okRight {
		return 0, 0, &RuntimeError{
			Token:   operator,
			Message: fmt.Sprintf("Operands must be numbers for operator '%s'.", operator.lexeme),
		}
	}

	return leftNumber, rightNumber, nil
}

func (i *Interpreter) stringify(obj any) string {
	if obj == nil {
		return "nil"
	}

	if num, ok := obj.(float64); ok {
		text := fmt.Sprintf("%v", num)
		if strings.HasSuffix(text, ".0") {
			text = text[:len(text)-2]
		}
		return text
	}
	return fmt.Sprintf("%v", obj)
}
