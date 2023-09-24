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
	receivedInput transaction.RecordRequest
	create        func(ctx context.Context, input transaction.RecordRequest) (string, error)
	receivedID    string
	get           func(ctx context.Context, id string) (*transaction.Retrieve, error)
}

func (s *stubService) Create(ctx context.Context, input transaction.RecordRequest) (string, error) {
	s.receivedInput = input
	return s.create(ctx, input)
}

func (s *stubService) Get(ctx context.Context, id string) (*transaction.Retrieve, error) {
	s.receivedID = id
	return s.get(ctx, id)
}

func TestTransaction_Store(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"

	mockSvc := &stubService{
		create: func(ctx context.Context, input transaction.RecordRequest) (string, error) {
			return id, nil
		},
	}

	input := transaction.RecordRequest{
		Description: "food",
		Amount:      23.12,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h := NewHandler(mockSvc)
	h.Store(w, req)

	var got transaction.RecordResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	want := transaction.RecordResponse{
		ID: id,
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, want, got)
	assert.Equal(t, input, mockSvc.receivedInput)
}

func TestTransaction_Store_Error(t *testing.T) {
	someErr := errors.New("some error")

	testCases := map[string]struct {
		reqBody        func() []byte
		mockSvc        *stubService
		wantStatusCode int
	}{
		"invalid json body": {
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
						Description: "food",
						Amount:      23.12,
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

	mockSvc := &stubService{
		get: func(ctx context.Context, id string) (*transaction.Retrieve, error) {
			return &transaction.Retrieve{
				ID:              id,
				Description:     "food",
				TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
				Amount:          23.12,
			}, nil
		},
	}

	path := fmt.Sprintf("/transactions/%s", id)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()

	h := NewHandler(mockSvc)
	r := chi.NewRouter()
	r.HandleFunc("/transactions/{id}", h.Retrieve)
	r.ServeHTTP(w, req)

	var got transaction.Retrieve
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	want := transaction.Retrieve{
		ID:              id,
		Description:     "food",
		TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
		Amount:          23.12,
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, want, got)
	assert.Equal(t, id, mockSvc.receivedID)
}

func TestTransaction_Retrieve_Error(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"
	someErr := errors.New("some error")

	testCases := map[string]struct {
		mockSvc        *stubService
		wantStatusCode int
	}{
		"validation error": {
			mockSvc: &stubService{
				get: func(ctx context.Context, id string) (*transaction.Retrieve, error) {
					return nil, httpresponse.ErrValidation
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		"service error": {
			mockSvc: &stubService{
				get: func(ctx context.Context, id string) (*transaction.Retrieve, error) {
					return nil, someErr
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		"not found error": {
			mockSvc: &stubService{
				get: func(ctx context.Context, id string) (*transaction.Retrieve, error) {
					return nil, httpresponse.ErrNotFound
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			path := fmt.Sprintf("/transactions/%s", id)
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
