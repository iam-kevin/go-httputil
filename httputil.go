package httputil

import (
	"encoding/json"
	"iam-kevin/errors"
	"log/slog"
	"net/http"
)

func ErrorWithStatus(w http.ResponseWriter, statusCode int, err interface{}) {
	slog.Error("failed", "status", statusCode, "error", err)
	var err_ error
	if e, ok := err.(error); ok {
		err_ = e
	} else {
		err_ = errors.New(err.(string))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      false,
		"message": err_.Error(),
	})
}

func Error(w http.ResponseWriter, err interface{}) {
	ErrorWithStatus(w, http.StatusInternalServerError, err)
}

func JsonWithStatus(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Json encodes data to response body with status 200 OK
func Json(w http.ResponseWriter, data interface{}) {
	JsonWithStatus(w, http.StatusOK, data)
}

// MessageWithStatus encodes message to response json with status 200 OK
func MessageWithStatus(w http.ResponseWriter, statusCode int, message string) {
	JsonWithStatus(w, statusCode, map[string]interface{}{
		"ok":      true,
		"message": message,
	})
}

// MessageWithStatus encodes message to response json with status 200 OK
func Message(w http.ResponseWriter, message string) {
	JsonWithStatus(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": message,
	})
}

// OK encodes an OK message to the response body with status 200 OK
func OK(w http.ResponseWriter) {
	JsonWithStatus(w, http.StatusOK, map[string]bool{
		"ok": true,
	})
}

func InternalError(w http.ResponseWriter, err error) {
	InternalErrorWithStatus(w, http.StatusInternalServerError, err)
}

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
