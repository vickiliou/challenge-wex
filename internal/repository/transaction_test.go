package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vickiliou/challenge-wex/internal/transaction"
)

func TestTransaction_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	txn := transaction.Transaction{
		ID:              "b62a64c9-0008-4148-99f6-9c8086a1dd42",
		Description:     "food",
		TransactionDate: time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC),
		Amount:          20.20,
	}

	mock.ExpectExec(`INSERT INTO "transaction" (id, description, date, amount)  VALUES (?, ?, ?, ?)`).
		WithArgs(txn.ID, txn.Description, txn.TransactionDate, txn.Amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(db)

	got, gotErr := repo.Create(context.Background(), txn)
	assert.NoError(t, gotErr)
	assert.Equal(t, txn.ID, got)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTransaction_Create_Error(t *testing.T) {
	wantErr := errors.New("some error")

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	txn := transaction.Transaction{
		ID:              "b62a64c9-0008-4148-99f6-9c8086a1dd42",
		Description:     "food",
		TransactionDate: time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC),
		Amount:          20.20,
	}

	mock.ExpectExec(`INSERT INTO "transaction"`).
		WithArgs(txn.ID, txn.Description, txn.TransactionDate, txn.Amount).
		WillReturnError(wantErr)

	repo := NewRepository(db)

	got, gotErr := repo.Create(context.Background(), txn)
	assert.Empty(t, got)
	assert.ErrorContains(t, gotErr, wantErr.Error())
}

func TestTransaction_FindByID(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	want := &transaction.Transaction{
		ID:              id,
		Description:     "food",
		TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
		Amount:          20.20,
	}

	row := mock.NewRows([]string{"id", "description", "date", "amount"}).
		AddRow(want.ID, want.Description, want.TransactionDate, want.Amount)

	mock.ExpectQuery(`SELECT id, description, date, amount FROM "transaction" WHERE id = ?`).
		WithArgs(id).
		WillReturnRows(row)

	repo := NewRepository(db)

	got, gotErr := repo.FindByID(context.Background(), id)
	assert.NoError(t, gotErr)
	assert.Equal(t, want, got)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTransaction_FindByID_Error(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	testCases := map[string]struct {
		rows    *sqlmock.Rows
		rowErr  error
		wantErr string
	}{
		"no rows error": {
			rowErr:  sql.ErrNoRows,
			rows:    mock.NewRows([]string{"description", "date", "amount"}),
			wantErr: "not found",
		},
		"row scan error": {
			rows:    mock.NewRows([]string{""}).AddRow(1),
			wantErr: "Scan",
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			mock.ExpectQuery(`SELECT (.+)`).
				WithArgs(id).
				WillReturnRows(tc.rows).
				WillReturnError(tc.rowErr)

			repo := NewRepository(db)

			got, gotErr := repo.FindByID(context.Background(), id)
			assert.Nil(t, got)
			assert.ErrorContains(t, gotErr, tc.wantErr)
		})
	}
}
