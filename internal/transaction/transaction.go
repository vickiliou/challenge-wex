package transaction

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

// Transaction represents a transaction stored in the database.
type Transaction struct {
	ID              string
	Description     string
	TransactionDate time.Time
	Amount          float64
}

// RecordRequest represents input data for a transaction request provided by the user.
type RecordRequest struct {
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
	Amount          float64   `json:"amount"`
}

// RecordResponse represents the response for a transaction request.
type RecordResponse struct {
	ID string `json:"id"`
}

// RetrieveResponse represents user transaction data.
type RetrieveResponse struct {
	ID              string    `json:"id"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
	OriginalAmount  float64   `json:"original_amount"`
	ExchangeRate    float64   `json:"exchange_rate"`
	ConvertedAmount float64   `json:"converted_amount"`
}

// validate checks if the request data is valid.
func (r *RecordRequest) validate() error {
	if len(r.Description) > 50 {
		return errors.New("description must not exceed 50 characters")
	}

	if !isValidDate(r.TransactionDate) {
		return errors.New("invalid date format")
	}

	if r.Amount < 0 {
		return errors.New("amount must be a positive number")
	}

	if math.Mod(r.Amount*100, 1) != 0 {
		return errors.New("amount must be rounded to the nearest cent")
	}
	return nil
}

// isValidUUID checks if a given string is a valid UUID.
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// isValidDate checks if date is a valid RFC3339 formatted timestamp.
func isValidDate(date time.Time) bool {
	_, err := time.Parse(time.RFC3339, date.Format(time.RFC3339))
	return err == nil
}

// roundTwoDecimal rounds a float number to two decimal places.
func roundTwoDecimal(amount float64) float64 {
	return math.Round(amount*100) / 100
}
