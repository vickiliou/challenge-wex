package transaction

import (
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
	testCases := map[string]struct {
		input   *RecordRequest
		wantErr string
	}{
		"description more than 50 characters": {
			input: &RecordRequest{
				Description: "more than 50 characters, more than 50 characters!!!",
			},
			wantErr: "description must not exceed 50 characters",
		},

		"invalid date": {
			input: &RecordRequest{
				TransactionDate: time.Date(-999, -99, -99, 0, 0, 0, 0, time.UTC),
			},
			wantErr: "invalid date format",
		},
		"negative amount": {
			input: &RecordRequest{
				Amount: -1,
			},
			wantErr: "amount must be a positive number",
		},
		"not rounded to the nearest cent": {
			input: &RecordRequest{
				Amount: 9.555,
			},
			wantErr: "amount must be rounded to the nearest cent",
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			gotErr := tc.input.validate()
			assert.EqualError(t, gotErr, tc.wantErr)
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
