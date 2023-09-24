package transaction

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

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

// Record represents a transaction record stored in the database.
type Record struct {
	ID              string
	Description     string
	TransactionDate time.Time
	Amount          float64
}

// Retrieve represents a retrieved transaction from the database.
type Retrieve struct {
	ID              string    `json:"id"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
	Amount          float64   `json:"amount"`
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

// validate checks if the store data is valid.
func (r *RecordRequest) validate() error {
	if len(r.Description) > 50 {
		return errors.New("description must not exceed 50 characters")
	}

	if !isValidDateFormat(r.TransactionDate) {
		return errors.New("invalid date")
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

func isValidDateFormat(date time.Time) bool {
	expectedFormat := "2006-01-02"
	formattedDate := date.Format(expectedFormat)
	return formattedDate != date.String()
}
