package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vickiliou/challenge-wex/internal/httpresponse"
	"github.com/vickiliou/challenge-wex/internal/transaction"
)

// Repository handles database operations for transactions.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository with the provided database connection.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create inserts a transaction record into the database.
func (r *Repository) Create(ctx context.Context, txn transaction.Transaction) (string, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO "transaction" 
			(id, description, date, amount) 
		VALUES 
			(?, ?, ?, ?)`,
		txn.ID, txn.Description, txn.TransactionDate, txn.Amount)

	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	return txn.ID, nil
}

// FindByID retrieves a transaction record by its ID from the database.
func (r *Repository) FindByID(ctx context.Context, id string) (*transaction.Transaction, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			id, description, date, amount
		FROM 
			"transaction" 
		WHERE 
			id = ?`,
		id)

	var txn transaction.Transaction
	if err := row.Scan(&txn.ID, &txn.Description, &txn.TransactionDate, &txn.Amount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w transaction ID %s", httpresponse.ErrNotFound, id)
		}
		return nil, fmt.Errorf("failed to retrieve transaction: %w", err)
	}

	return &txn, nil
}
