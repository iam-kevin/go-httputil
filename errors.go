package httputil

// httperror implements the HttpError interface and represents an HTTP error
// with a status code and underlying error.
type httperror struct {
	status int
	err    error
}

// Status returns the HTTP status code associated with this error.
func (he httperror) Status() int {
	return he.status
}

// Error returns the error message string.
// This implements the standard error interface.
func (he httperror) Error() string {
	return he.err.Error()
}

// Cause returns the underlying error that caused this HTTP error.
// Returns nil if the underlying error doesn't implement ErrorWithCause.
func (he httperror) Cause() error {
	// return errors.Unwarp(he.err)
	return he.err
}

// NewError creates a new HTTP error with the specified status code and underlying error.
// The returned error implements the HttpError interface.
//
// Example:
//
//	err := NewError(404, errors.New("user not found"))
//	if httpErr, ok := err.(HttpError); ok {
//		fmt.Printf("Status: %d, Message: %s", httpErr.Status(), httpErr.Error())
//	}
func NewError(status int, err error) error {
	return &httperror{
		status: status,
		err:    err,
	}
}
