package httputil

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/iam-kevin/go-errors"
)

// ErrorWithStatus sends a JSON error response with the specified HTTP status code.
// The error can be either a string or an error type.
//
// The response format is:
//
//	{
//		"ok": false,
//		"message": "error message"
//	}
//
// Example:
//
//	httputil.ErrorWithStatus(w, 400, "invalid input")
//	httputil.ErrorWithStatus(w, 404, errors.New("user not found"))
func ErrorWithStatus(w http.ResponseWriter, statusCode int, err interface{}) {
	slog.Error("failed", "status", statusCode, "error", err)
	var err_ error
	switch e := err.(type) {
	case error:
		err_ = e
	case string:
		err_ = errors.New(e)
	case httperror:
		err_ = e.Cause()
	default:
		err_ = errors.New("unknown error occured")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      false,
		"message": err_.Error(),
	})
}

// Error sends a JSON error response with HTTP 500 Internal Server Error status.
// This is a convenience function equivalent to ErrorWithStatus with status 500.
//
// Example:
//
//	httputil.Error(w, "something went wrong")
func Error(w http.ResponseWriter, err interface{}) {
	ErrorWithStatus(w, http.StatusInternalServerError, err)
}

// JsonWithStatus encodes data as JSON and sends it with the specified HTTP status code.
// Sets appropriate Content-Type headers and security headers.
//
// Example:
//
//	user := User{Name: "John", Email: "john@example.com"}
//	httputil.JsonWithStatus(w, 201, user)
func JsonWithStatus(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Json encodes data as JSON and sends it with HTTP 200 OK status.
// This is a convenience function equivalent to JsonWithStatus with status 200.
//
// Example:
//
//	users := []User{{Name: "John"}, {Name: "Jane"}}
//	httputil.Json(w, users)
func Json(w http.ResponseWriter, data interface{}) {
	JsonWithStatus(w, http.StatusOK, data)
}

// MessageWithStatus sends a JSON message response with the specified HTTP status code.
// The response format is:
//
//	{
//		"ok": true,
//		"message": "your message"
//	}
//
// Example:
//
//	httputil.MessageWithStatus(w, 201, "user created successfully")
func MessageWithStatus(w http.ResponseWriter, statusCode int, message string) {
	JsonWithStatus(w, statusCode, map[string]interface{}{
		"ok":      true,
		"message": message,
	})
}

// Message sends a JSON message response with HTTP 200 OK status.
// This is a convenience function equivalent to MessageWithStatus with status 200.
// The response format is:
//
//	{
//		"ok": true,
//		"message": "your message"
//	}
//
// Example:
//
//	httputil.Message(w, "operation completed")
func Message(w http.ResponseWriter, message string) {
	JsonWithStatus(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": message,
	})
}

// OK sends a simple success response with HTTP 200 OK status.
// The response format is:
//
//	{
//		"ok": true
//	}
//
// Example:
//
//	httputil.OK(w)
func OK(w http.ResponseWriter) {
	JsonWithStatus(w, http.StatusOK, map[string]bool{
		"ok": true,
	})
}

// InternalError sends an internal server error response with HTTP 500 status.
// The error details are logged but not exposed to the client for security reasons.
// This is a convenience function equivalent to InternalErrorWithStatus with status 500.
//
// Example:
//
//	err := database.Connect()
//	if err != nil {
//		httputil.InternalError(w, err)
//		return
//	}
func InternalError(w http.ResponseWriter, err error) {
	InternalErrorWithStatus(w, http.StatusInternalServerError, err)
}

// InternalErrorWithStatus sends an internal server error response with the specified HTTP status code.
// The error details are logged with full context but not exposed to the client for security reasons.
// The client receives a generic "internal server error" message.
//
// If the error implements ErrorWithCause, the underlying cause is also logged.
//
// The response format is:
//
//	{
//		"ok": false,
//		"message": "internal server error"
//	}
//
// Example:
//
//	err := criticalSystemFailure()
//	if err != nil {
//		httputil.InternalErrorWithStatus(w, 503, err)
//		return
//	}
func InternalErrorWithStatus(w http.ResponseWriter, status int, err error) {
	if errwc, ok := err.(errors.ErrorWithCause); ok {
		slog.Error("internal error: "+errwc.Error(), "cause", errwc.Cause())
	} else {
		slog.Error("internal error: " + err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      false,
		"message": "internal server error",
	})
}
