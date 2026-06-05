package stoplang

// RuntimeError represents the data layout of a runtime evaluation crash.
type RuntimeError struct {
	Token   Token
	Message string
}

// Error implements the built-in Go error interface.
func (e *RuntimeError) Error() string {
	return e.Message
}
