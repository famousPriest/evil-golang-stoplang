package stoplang

import (
	"fmt"
	"os"
	"strings"
)

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(statements []Stmt, lang *Lang) {
	for _, statement := range statements {
		err := i.execute(statement, lang)
		if err != nil {
			if runtimeErr, ok := err.(*RuntimeError); ok {
				fmt.Fprintf(os.Stderr, "%s\n[line %d]\n", runtimeErr.Message, runtimeErr.Token.line)
				lang.hadRuntimeError = true
			} else {
				fmt.Fprintf(os.Stderr, "Unknown error: %v\n", err)
			}
			return
		}
	}
}

func (i *Interpreter) execute(stmt Stmt, lang *Lang) error {
	_, err := stmt.Accept(i)
	return err
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
				return nil, &RuntimeError{
					Token:   e.Operator,
					Message: "Operand must be a number.",
				}
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
				Message: "Operands must be two numbers or two strings.",
			}

		case BANG_EQUAL:
			return !i.isEqual(left, right), nil

		case EQUAL_EQUAL:
			return i.isEqual(left, right), nil
		}

		return nil, fmt.Errorf("unknown binary operator: %v", e.Operator)

	case *Assign:
		value, err := i.Evaluate(e.Value)
		if err != nil {
			return nil, err
		}

		err = i.environment.Assign(e.Name, value)
		if err != nil {
			return nil, err
		}

		return value, nil

	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (i *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) (any, error) {
	_, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) (any, error) {
	for {
		condition, err := i.Evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !i.isTruthy(condition) {
			break
		}
		return i.execute(stmt.Body, nil), err
	}
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *IfStmt) (any, error) {
	condition, err := i.Evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(condition) {
		return i.execute(stmt.ThenBranch, nil), err
	}
	if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch, nil), err
	}
	return nil, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) (any, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.tokenType == OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	return right, nil
}

func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt) (any, error) {
	value, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(i.stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(stmt *VarStmt) (any, error) {
	var value any
	var err error

	if stmt.Initializer != nil {
		value, err = i.Evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.environment.Define(stmt.Name.lexeme, value)

	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) (any, error) {
	scopedEnv := NewEnclosingEnvironment(i.environment)

	return nil, i.executeBlock(stmt.Statements, scopedEnv)
}

func (i *Interpreter) executeBlock(statements []Stmt, env *Environment) error {
	previousEnv := i.environment

	defer func() {
		i.environment = previousEnv
	}()

	i.environment = env

	for _, statement := range statements {
		err := i.execute(statement, nil)
		if err != nil {
			return err
		}
	}

	return nil
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
