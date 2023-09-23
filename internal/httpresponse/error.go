package httpresponse

import (
	"errors"

	"golang.org/x/exp/slog"
)

var (
	// ErrValidation is an error indicating a validation failure.
	ErrValidation = errors.New("validation error")

	// ErrNotFound is an error indicating that a resource was not found.
	ErrNotFound = errors.New("not found")
)

// LogError logs an error with additional information.
func LogError(msg string, statusCode int, err error) {
	slog.Error(
		msg,
		slog.Int("status_code", statusCode),
		slog.String("error", err.Error()),
	)
}
