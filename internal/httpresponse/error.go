package httpresponse

import (
	"errors"

	"golang.org/x/exp/slog"
)

var (
	// ErrValidation indicates a validation failure.
	ErrValidation = errors.New("validation error")

	// ErrNotFound indicates that a resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrNoCurrencyConversion indicates that no currency conversion rate data was available within 6 months before the purchase date.
	ErrNoCurrencyConversion = errors.New("no currency conversion rate available within 6 months before the purchase date")

	// ErrInvalidRequestPayload indicates that the http request payload is invalid.
	ErrInvalidRequestPayload = errors.New("invalid request payload")
)

// LogError logs an error with additional information.
func LogError(msg string, statusCode int, err error) {
	slog.Error(
		msg,
		slog.Int("status_code", statusCode),
		slog.String("error", err.Error()),
	)
}
