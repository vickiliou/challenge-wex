package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Validate(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		input := &Request{
			Description: "food",
			Amount:      40.50,
		}
		gotErr := input.validate()
		assert.Nil(t, gotErr)
	})
}

func TestTransaction_Validate_Error(t *testing.T) {
	testCases := map[string]struct {
		input   *Request
		wantErr string
	}{
		"description more than 50 characters": {
			input: &Request{
				Description: "more than 50 characters, more than 50 characters!!!",
			},
			wantErr: "description must not exceed 50 characters",
		},
		"negative amount": {
			input: &Request{
				Amount: -1,
			},
			wantErr: "amount must be a positive number",
		},
		"not rounded to the nearest cent": {
			input: &Request{
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
