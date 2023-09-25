package transaction

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Validate(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		input := &RecordRequest{
			Description:     "food",
			TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
			Amount:          40.50,
		}
		gotErr := input.validate()
		assert.Nil(t, gotErr)
	})
}

func TestTransaction_Validate_Error(t *testing.T) {
	description := "food"
	transactionDate := time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC)
	amount := 40.50

	testCases := map[string]struct {
		input   *RecordRequest
		wantErr string
	}{
		"empty transaction": {
			input:   nil,
			wantErr: "required",
		},
		"empty description": {
			input: &RecordRequest{
				Description:     "",
				TransactionDate: transactionDate,
				Amount:          amount,
			},
			wantErr: "required",
		},
		"description more than 50 characters": {
			input: &RecordRequest{
				Description: "more than 50 characters, more than 50 characters!!!",
			},
			wantErr: "not exceed 50 characters",
		},
		"empty date": {
			input: &RecordRequest{
				Description:     description,
				TransactionDate: time.Time{},
				Amount:          amount,
			},
			wantErr: "required",
		},
		"invalid date": {
			input: &RecordRequest{
				Description:     description,
				TransactionDate: time.Date(-999, -99, -99, 0, 0, 0, 0, time.UTC),
				Amount:          40.50,
			},
			wantErr: "invalid date format",
		},
		"empty amount": {
			input: &RecordRequest{
				Description:     description,
				TransactionDate: transactionDate,
				Amount:          math.NaN(),
			},
			wantErr: "required",
		},
		"negative amount": {
			input: &RecordRequest{
				Description:     description,
				TransactionDate: transactionDate,
				Amount:          -1,
			},
			wantErr: "positive number",
		},
		"not rounded to two decimal places": {
			input: &RecordRequest{
				Description:     description,
				TransactionDate: transactionDate,
				Amount:          9.5579,
			},
			wantErr: "two decimal places",
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			gotErr := tc.input.validate()
			assert.ErrorContains(t, gotErr, tc.wantErr)
		})
	}
}

func TestTransaction_IsValidUUID(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  bool
	}{
		"valid UUID": {
			input: "b62a64c9-0008-4148-99f6-9c8086a1dd42",
			want:  true,
		},
		"invalid UUID": {
			input: "invalid-uuid",
			want:  false,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			got := isValidUUID(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}

}
