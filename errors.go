package httputil

import "github.com/iam-kevin/go-errors"

type httperror struct {
	status int
	err    error
}

func (he httperror) Status() int {
	return he.status
}

func (he httperror) Error() string {
	return he.err.Error()
}

func (he httperror) Cause() error {
	if ewc, ok := he.err.(errors.ErrorWithCause); ok {
		return ewc
	}

	return nil
}

func NewError(status int, err error) error {
	return &httperror{
		status: status,
		err:    err,
	}
}
