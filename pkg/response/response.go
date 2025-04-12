package response

import (
	"errors"
)

type Error struct {
	Code int
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NewError(code int, err string) error {
	return &Error{code, errors.New(err)}
}
