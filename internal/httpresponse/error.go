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

	// // ErrConvertTargetCurrency is returned when a currency cannot be converted to the target currency due to an empty exchange rate.
	ErrConvertTargetCurrency = errors.New("cannot be converted to the target currency")

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
