package internal

import "fmt"

type Error struct {
	Cause   error
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func NewError(cause error, format string, args ...any) error {
	return &Error{
		Cause:   cause,
		Message: fmt.Sprintf(format, args...),
	}
}
