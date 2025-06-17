package httputil

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/iam-kevin/go-errors"
)

// toErr converts various error types to a standard error interface.
// It accepts string, error, or any other type and converts them to an error.
func toErr(err interface{}) error {
	var er error
	switch v := err.(type) {
	case string:
		{
			er = errors.New(v)
		}
	case error:
		{
			er = v
		}
	default:
		{
			er = errors.New("missing error details")
		}
	}

	return er
}

// Assert performs an assertion within a request handler.
// If the condition is false, it panics with an HTTP 500 Internal Server Error.
//
// This assertion must be captured by the MiddlewareHTTPAssertionRecoverer middleware
// to properly format the error response.
//
// Example:
//
//	user := getUserFromDB(userID)
//	httputil.Assert(user != nil, "user not found")
func Assert(condition bool, err interface{}) {
	AssertWithStatus(condition, http.StatusInternalServerError, err)
}

// AssertWithStatus performs an assertion within a request handler with a custom HTTP status code.
// If the condition is false, it panics with the specified HTTP status code.
//
// This assertion must be captured by the MiddlewareHTTPAssertionRecoverer middleware
// to properly format the error response.
//
// Example:
//
//	httputil.AssertWithStatus(user.IsActive, 403, "user account is disabled")
func AssertWithStatus(condition bool, status int, err interface{}) {
	if !condition {
		log.Printf("AssertionError(HTTP: %v): %s", status, err)
		panic(httperror{
			status: status,
			err:    toErr(err),
		})
	}
}

// AssertErrorIsNilWithStatus asserts that err == nil. If not, panics with the specified HTTP status code.
//
// This assertion must be captured by the MiddlewareHTTPAssertionRecoverer middleware
// to properly format the error response.
//
// Example:
//
//	err := validateInput(data)
//	httputil.AssertErrorIsNilWithStatus(400, err)
func AssertErrorIsNilWithStatus(status int, err interface{}) {
	if err != nil {
		log.Printf("AssertionError(HTTP: %v): %s", status, err)
		panic(httperror{
			status: status,
			err:    toErr(err),
		})
	}
}

// AssertErrorIsNil asserts that err == nil, otherwise panics with HTTP 500 Internal Server Error.
//
// This is a convenience function equivalent to AssertErrorIsNilWithStatus with status 500.
//
// Example:
//
//	err := processData(input)
//	httputil.AssertErrorIsNil(err)
func AssertErrorIsNil(err error) {
	AssertErrorIsNilWithStatus(http.StatusInternalServerError, err)
}

// HttpError represents an HTTP error with additional context.
// It provides the HTTP status code, error message, and optional underlying cause.
type HttpError interface {
	// Status returns the HTTP status code for this error
	Status() int
	// Error returns the error message
	Error() string
	// Cause returns the underlying error that caused this HTTP error, if any
	Cause() error
}

// MiddlewareHTTPAssertionRecoverer is a middleware that intercepts HTTP assertion panics
// and converts them into proper HTTP error responses.
//
// It catches panics from Assert, AssertWithStatus, and AssertErrorIsNil functions,
// and returns appropriate JSON error responses with the correct status codes.
//
// Status codes >= 500 are treated as internal errors and logged with full details,
// while client errors (< 500) are returned with the original error message.
//
// Example:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/api/users", userHandler)
//
//	server := &http.Server{
//		Handler: httputil.MiddlewareHTTPAssertionRecoverer(mux),
//	}
func MiddlewareHTTPAssertionRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, cancel := context.WithCancel(r.Context())
		defer func() {
			if r := recover(); r != nil {
				if herr, ok := r.(HttpError); ok {
					if herr.Status() >= http.StatusInternalServerError {
						InternalErrorWithStatus(w, herr.Status(), herr)
					} else {
						ErrorWithStatus(w, herr.Status(), herr)
					}
					cancel()
				} else {
					slog.Error("unexpected panic", "error", r)
					panic(r)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
