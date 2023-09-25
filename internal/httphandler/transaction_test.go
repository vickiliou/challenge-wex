package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vickiliou/challenge-wex/internal/httpresponse"
	"github.com/vickiliou/challenge-wex/internal/transaction"
)

type stubService struct {
	receivedRecordRequest    transaction.RecordRequest
	create                   func(ctx context.Context, input transaction.RecordRequest) (string, error)
	receivedRetrievedRequest transaction.RetrieveRequest
	get                      func(ctx context.Context, input transaction.RetrieveRequest) (*transaction.RetrieveResponse, error)
}

func (s *stubService) Create(ctx context.Context, input transaction.RecordRequest) (string, error) {
	s.receivedRecordRequest = input
	return s.create(ctx, input)
}

func (s *stubService) Get(ctx context.Context, input transaction.RetrieveRequest) (*transaction.RetrieveResponse, error) {
	s.receivedRetrievedRequest = input
	return s.get(ctx, input)
}

func TestTransaction_Store(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"

	mockSvc := &stubService{
		create: func(ctx context.Context, input transaction.RecordRequest) (string, error) {
			return id, nil
		},
	}

	input := transaction.RecordRequest{
		Description:     "food",
		TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
		Amount:          23.12,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h := NewHandler(mockSvc)
	h.Store(w, req)

	var got transaction.RecordResponse
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)

	want := transaction.RecordResponse{
		ID: id,
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, want, got)
	assert.Equal(t, input, mockSvc.receivedRecordRequest)
}

func TestTransaction_Store_Error(t *testing.T) {
	someErr := errors.New("some error")

	testCases := map[string]struct {
		reqBody        func() []byte
		mockSvc        *stubService
		wantStatusCode int
	}{
		"invalid json request body": {
			reqBody: func() []byte {
				jsonValue, _ := json.Marshal(",")
				return jsonValue
			},
			mockSvc: &stubService{
				create: func(ctx context.Context, input transaction.RecordRequest) (string, error) {
					return "", nil
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		"validation error": {
			reqBody: func() []byte {
				jsonValue, _ := json.Marshal(
					transaction.RecordRequest{
						Description: "food",
						Amount:      23.12,
					})
				return jsonValue
			},
			mockSvc: &stubService{
				create: func(ctx context.Context, input transaction.RecordRequest) (string, error) {
					return "", httpresponse.ErrValidation
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		"service error": {
			reqBody: func() []byte {
				jsonValue, _ := json.Marshal(
					transaction.RecordRequest{
						Description:     "food",
						TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
						Amount:          23.12,
					})
				return jsonValue
			},
			mockSvc: &stubService{
				create: func(ctx context.Context, input transaction.RecordRequest) (string, error) {
					return "", someErr
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(tc.reqBody()))
			w := httptest.NewRecorder()

			h := NewHandler(tc.mockSvc)
			h.Store(w, req)

			assert.Equal(t, tc.wantStatusCode, w.Code)
		})
	}
}

func TestTransaction_Retrieve(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"
	countryCurrency := "Brazil"
	currency := "Real"

	want := transaction.RetrieveResponse{
		ID:              id,
		Description:     "food",
		TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
		OriginalAmount:  23.12,
		ExchangeRate:    3.456,
		ConvertedAmount: 79.90,
	}

	mockSvc := &stubService{
		get: func(ctx context.Context, input transaction.RetrieveRequest) (*transaction.RetrieveResponse, error) {
			return &want, nil
		},
	}

	path := fmt.Sprintf("/transactions/%s?country_currency=%s&currency=%s", id, countryCurrency, currency)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()

	h := NewHandler(mockSvc)
	r := chi.NewRouter()
	r.HandleFunc("/transactions/{id}", h.Retrieve)
	r.ServeHTTP(w, req)

	var got transaction.RetrieveResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	wantRetrievedRequest := transaction.RetrieveRequest{
		ID:              id,
		CountryCurrency: countryCurrency,
		Currency:        currency,
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, want, got)
	assert.Equal(t, wantRetrievedRequest, mockSvc.receivedRetrievedRequest)
}

func TestTransaction_Retrieve_Error(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"
	countryCurrency := "Brazil"
	currency := "Real"
	someErr := errors.New("some error")

	testCases := map[string]struct {
		mockSvc        *stubService
		wantStatusCode int
	}{
		"validation error": {
			mockSvc: &stubService{
				get: func(ctx context.Context, input transaction.RetrieveRequest) (*transaction.RetrieveResponse, error) {
					return nil, httpresponse.ErrValidation
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		"not found error": {
			mockSvc: &stubService{
				get: func(ctx context.Context, input transaction.RetrieveRequest) (*transaction.RetrieveResponse, error) {
					return nil, httpresponse.ErrNotFound
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
		"no exchange rate": {
			mockSvc: &stubService{
				get: func(ctx context.Context, input transaction.RetrieveRequest) (*transaction.RetrieveResponse, error) {
					return nil, httpresponse.ErrConvertTargetCurrency
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		"service error": {
			mockSvc: &stubService{
				get: func(ctx context.Context, input transaction.RetrieveRequest) (*transaction.RetrieveResponse, error) {
					return nil, someErr
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			path := fmt.Sprintf("/transactions/%s?country_currency=%s&currency=%s", id, countryCurrency, currency)

			req := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()

			h := NewHandler(tc.mockSvc)
			r := chi.NewRouter()
			r.HandleFunc("/transactions/{id}", h.Retrieve)
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatusCode, w.Code)
		})
	}
}
