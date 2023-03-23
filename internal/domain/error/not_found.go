package errorx

import "errors"

type NotFound struct {
	err error
}

// NewNotFound returns new NotFound instance.
func NewNotFound() *NotFound {
	return &NotFound{}
}

// WrapInNotFound wraps err into NotFound.
func WrapInNotFound(err error) *NotFound {
	return &NotFound{
		err: err,
	}
}

// IsNotFound checks if err is instance of NotFound.
func IsNotFound(err error) bool {
	var e *NotFound
	return errors.As(err, &e)
}

func (e *NotFound) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "not found"
}

// Unwrap unwraps error. See errors.Unwrap for more details.
func (e *NotFound) Unwrap() error {
	return e.err
}
