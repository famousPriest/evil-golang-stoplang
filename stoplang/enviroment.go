package stoplang

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    make(map[string]any),
	}
}

func NewEnclosingEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) Get(name Token) (any, error) {
	if value, ok := e.values[name.lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, &RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("Undefined variable '%s'.", name.lexeme),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name Token, value any) error {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	return &RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("Undefined variable '%s'.", name.lexeme),
	}
}
