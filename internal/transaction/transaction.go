package transaction

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

// Request represents a transaction request provided by the user.
type Request struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

// RequestResponse represents a response for a transaction request.
type RequestResponse struct {
	ID string `json:"id"`
}

// Record represents a transaction record stored in the database.
type Record struct {
	ID          string
	Description string
	Amount      float64
}

// Retrieve represents a retrieved transaction.
type Retrieve struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Amount      float64   `json:"amount"`
}

// validate checks if the store data is valid.
func (r *Request) validate() error {
	if len(r.Description) > 50 {
		return errors.New("description must not exceed 50 characters")
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
