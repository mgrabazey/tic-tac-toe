package errorx

import "errors"

type BadRequest struct {
	err error
}

// NewBadRequest returns new BadRequest instance.
func NewBadRequest() *BadRequest {
	return &BadRequest{}
}

// WrapInBadRequest wraps err into BadRequest.
func WrapInBadRequest(err error) *BadRequest {
	return &BadRequest{
		err: err,
	}
}

// IsBadRequest checks if err is instance of BadRequest.
func IsBadRequest(err error) bool {
	var e *BadRequest
	return errors.As(err, &e)
}

func (e *BadRequest) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "bad request"
}

// Unwrap unwraps error. See errors.Unwrap for more details.
func (e *BadRequest) Unwrap() error {
	return e.err
}
