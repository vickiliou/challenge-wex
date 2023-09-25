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

// RetrieveRequest represents a request to retrieve user transaction data.
type RetrieveRequest struct {
	ID              string `json:"id"`
	CurrencyCountry string `json:"currency_country"`
	CurrencyCode    string `json:"currency_code"`
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

// validate checks if the record request data is valid.
func (r *RecordRequest) validate() error {
	if r == nil {
		return errors.New("all fields are required")
	}

	if err := validateDescription(r.Description); err != nil {
		return err
	}

	if err := validateTransactionDate(r.TransactionDate); err != nil {
		return err
	}

	if err := validateAmount(r.Amount); err != nil {
		return err
	}

	return nil
}

// validate checks if the retrieve request data is valid.
func (r *RetrieveRequest) validate() error {
	if isValidUUID(r.ID) {
		return errors.New("invalid UUID")
	}
	if isEmpty(r.CurrencyCountry) {
		return errors.New("currency country is required")
	}

	if isEmpty(r.CurrencyCode) {
		return errors.New("currency code is required")
	}
	return nil
}

// validateDescription checks if the description field is a valid RFC3339 formatted timestamp and not empty.
func validateDescription(description string) error {
	if isEmpty(description) {
		return errors.New("description is required")
	}

	if len(description) > 50 {
		return errors.New("description must not exceed 50 characters")
	}

	return nil
}

// validateTransactionDate checks if the transaction date field is valid and not empty.
func validateTransactionDate(transactionDate time.Time) error {
	if transactionDate.IsZero() {
		return errors.New("transaction date is required")
	}

	if _, err := time.Parse(time.RFC3339, transactionDate.Format(time.RFC3339)); err != nil {
		return errors.New("invalid date format")
	}

	return nil
}

// validateAmount checks if the amount field is valid and not empty.
func validateAmount(amount float64) error {
	if math.IsNaN(amount) {
		return errors.New("amount is required")
	}

	if amount <= 0 {
		return errors.New("amount must be a positive number")
	}

	if amount != roundTwoDecimal(amount) {
		return errors.New("amount must be rounded to two decimal places")
	}

	return nil
}

// isValidUUID checks if a given string is a valid UUID.
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err != nil
}

// isEmpty checks if a given string is empty.
func isEmpty(s string) bool {
	return len(s) == 0
}

// roundTwoDecimal rounds a float number to two decimal places.
func roundTwoDecimal(amount float64) float64 {
	return math.Round(amount*100) / 100
}
