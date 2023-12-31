package transaction

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vickiliou/challenge-wex/internal/gateway"
	"github.com/vickiliou/challenge-wex/internal/httpresponse"
)

type stubRepository struct {
	receivedCreateInput Transactions
	create              func(ctx context.Context, txn Transactions) (string, error)
	receivedFindInput   string
	findByID            func(ctx context.Context, id string) (*Transactions, error)
}

func (s *stubRepository) Create(ctx context.Context, txn Transactions) (string, error) {
	s.receivedCreateInput = txn
	return s.create(ctx, txn)
}

func (s *stubRepository) FindByID(ctx context.Context, id string) (*Transactions, error) {
	s.receivedFindInput = id
	return s.findByID(ctx, id)
}

type stubGateway struct {
	receivedGwInput gateway.CurrencyExchangeRateRequest
	getExchangeRate func(input gateway.CurrencyExchangeRateRequest) (*gateway.CurrencyExchangeRate, error)
}

func (s *stubGateway) GetExchangeRate(input gateway.CurrencyExchangeRateRequest) (*gateway.CurrencyExchangeRate, error) {
	s.receivedGwInput = input
	return s.getExchangeRate(input)
}

func TestService_Create(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"

	mockRepo := &stubRepository{
		create: func(ctx context.Context, txn Transactions) (string, error) {
			return id, nil
		},
	}

	mockIDGen := func() string {
		return id
	}

	input := RecordRequest{
		Description:     "food",
		TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
		Amount:          20.47,
	}

	want := Transactions{
		ID:              id,
		Description:     input.Description,
		TransactionDate: input.TransactionDate,
		Amount:          input.Amount,
	}

	svc := NewService(mockRepo, nil, mockIDGen)
	got, gotErr := svc.Create(context.Background(), input)
	assert.NoError(t, gotErr)
	assert.Equal(t, id, got)
	assert.Equal(t, want, mockRepo.receivedCreateInput)
}

func TestService_Create_Error(t *testing.T) {
	someErr := errors.New("some error")

	testCases := map[string]struct {
		input    RecordRequest
		mockRepo *stubRepository
		wantErr  error
	}{
		"validation error": {
			input: RecordRequest{
				Description:     "food",
				TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
				Amount:          -5,
			},
			mockRepo: &stubRepository{
				create: func(ctx context.Context, txn Transactions) (string, error) {
					return "", nil
				},
			},
			wantErr: httpresponse.ErrValidation,
		},
		"repository error": {
			input: RecordRequest{
				Description:     "food",
				TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
				Amount:          20.47,
			},
			mockRepo: &stubRepository{
				create: func(ctx context.Context, txn Transactions) (string, error) {
					return "", someErr
				},
			},
			wantErr: someErr,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			mockIDGen := func() string {
				return "b62a64c9-0008-4148-99f6-9c8086a1dd42"
			}

			svc := NewService(tc.mockRepo, nil, mockIDGen)
			got, gotErr := svc.Create(context.Background(), tc.input)
			assert.Empty(t, got)
			assert.ErrorContains(t, gotErr, tc.wantErr.Error())
		})
	}
}

func TestService_Get(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"
	retrieve := &Transactions{
		ID:              id,
		Description:     "food",
		TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
		Amount:          23.12,
	}

	mockRepo := &stubRepository{
		findByID: func(ctx context.Context, id string) (*Transactions, error) {
			return retrieve, nil
		},
	}

	mockIDGen := func() string {
		return id
	}

	mockGw := &stubGateway{
		getExchangeRate: func(input gateway.CurrencyExchangeRateRequest) (*gateway.CurrencyExchangeRate, error) {
			return &gateway.CurrencyExchangeRate{
				CountryCurrencyDesc: "Brazil-Real",
				ExchangeRate:        "3.456",
			}, nil
		},
	}

	input := RetrieveRequest{
		ID:       id,
		Country:  "Brazil",
		Currency: "Real",
	}

	svc := NewService(mockRepo, mockGw, mockIDGen)
	got, gotErr := svc.Get(context.Background(), input)
	assert.NoError(t, gotErr)

	want := &RetrieveResponse{
		ID:              retrieve.ID,
		Description:     retrieve.Description,
		TransactionDate: retrieve.TransactionDate,
		OriginalAmount:  retrieve.Amount,
		ExchangeRate:    3.456,
		ConvertedAmount: 79.90,
	}

	wantGwInput := gateway.CurrencyExchangeRateRequest{
		TransactionDate: retrieve.TransactionDate,
		Country:         input.Country,
		Currency:        input.Currency,
	}

	assert.Equal(t, want, got)
	assert.Equal(t, id, mockRepo.receivedFindInput)
	assert.Equal(t, wantGwInput, mockGw.receivedGwInput)
}

func TestService_Get_Error(t *testing.T) {
	someErr := errors.New("some error")

	testCases := map[string]struct {
		input    RetrieveRequest
		mockRepo *stubRepository
		mockGw   *stubGateway
		wantErr  error
	}{
		"validation error": {
			input: RetrieveRequest{
				ID:       "invalid-uuid",
				Country:  "Brazil",
				Currency: "Real",
			},
			mockRepo: &stubRepository{
				findByID: func(ctx context.Context, id string) (*Transactions, error) {
					return nil, nil
				},
			},
			mockGw:  &stubGateway{},
			wantErr: httpresponse.ErrValidation,
		},
		"repository error": {
			input: RetrieveRequest{
				ID:       "b62a64c9-0008-4148-99f6-9c8086a1dd42",
				Country:  "Brazil",
				Currency: "Real",
			},
			mockRepo: &stubRepository{
				findByID: func(ctx context.Context, id string) (*Transactions, error) {
					return nil, someErr
				},
			},
			mockGw:  &stubGateway{},
			wantErr: someErr,
		},
		"gateway error": {
			input: RetrieveRequest{
				ID:       "b62a64c9-0008-4148-99f6-9c8086a1dd42",
				Country:  "Brazil",
				Currency: "Real",
			},
			mockRepo: &stubRepository{
				findByID: func(ctx context.Context, id string) (*Transactions, error) {
					return &Transactions{
						ID:              id,
						Description:     "food",
						TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
						Amount:          23.12,
					}, nil
				},
			},
			mockGw: &stubGateway{
				getExchangeRate: func(input gateway.CurrencyExchangeRateRequest) (*gateway.CurrencyExchangeRate, error) {
					return nil, someErr
				},
			},
			wantErr: someErr,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			mockIDGen := func() string {
				return "b62a64c9-0008-4148-99f6-9c8086a1dd42"
			}

			svc := NewService(tc.mockRepo, tc.mockGw, mockIDGen)
			got, gotErr := svc.Get(context.Background(), tc.input)
			assert.Nil(t, got)
			assert.ErrorContains(t, gotErr, tc.wantErr.Error())
		})
	}
}
