package httputil

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/iam-kevin/go-errors"
)

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

// Assert is used in performing an  assersion within a request handler.
// Status that's used with this assersion is 500 (Internal Server Error)
//
// This assersion must be captured by the `MiddlewareHTTPAssersionRecoverer` to
// ultimately decorates the errors to a format for the response body
func Assert(condition bool, err interface{}) {
	AssertWithStatus(condition, http.StatusInternalServerError, err)
}

// AssertWithStatus is used in performing an  assersion within a request handler.
//
// This assersion must be captured by the `MiddlewareHTTPAssersionRecoverer` to
// ultimately decorates the errors to a format for the response body
func AssertWithStatus(condition bool, status int, err interface{}) {
	if !condition {
		log.Printf("AssersionError(HTTP: %v): %s", status, err)
		panic(httperror{
			status: status,
			err:    toErr(err),
		})
	}
}

// AssertErrorIsNilWithStatus asserts the `err == nil`. If not, panics with status
//
// This assersion must be captured by the `MiddlewareHTTPAssersionRecoverer` to
// ultimately decorates the errors to a format for the response body
func AssertErrorIsNilWithStatus(status int, err interface{}) {
	if err != nil {
		log.Printf("AssersionError(HTTP: %v): %s", status, err)
		panic(httperror{
			status: status,
			err:    toErr(err),
		})
	}
}

// AssertErrorIsNil asserts that `err == nil`, otherwise panics with a status of 500
func AssertErrorIsNil(err error) {
	AssertErrorIsNilWithStatus(http.StatusInternalServerError, err)
}

// Error containing more information about the failed http request
type HttpError interface {
	Status() int
	Error() string
	Cause() error
}

// MiddlewareHTTPAssersionRecoverer is a middleware that intercepts a HTTP Assersion
// and return a response with the appropriate error and status
func MiddlewareHTTPAssersionRecoverer(next http.Handler) http.Handler {
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
					slog.Error("update", "error", r)
					panic(r)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
