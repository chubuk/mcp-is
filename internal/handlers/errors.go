package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("resource already exists")
	ErrValidation   = errors.New("validation failed")
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(code int, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		writeJSON(w, appErr.Code, appErr)
		return
	}
	if errors.Is(err, ErrNotFound) {
		writeJSON(w, http.StatusNotFound, AppError{Message: "not found"})
		return
	}
	slog.Error("unhandled error", "error", err)
	writeJSON(w, http.StatusInternalServerError, AppError{Message: "internal error"})
}
