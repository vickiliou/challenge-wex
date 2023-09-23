package transaction

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vickiliou/challenge-wex/internal/httpresponse"
)

type stubRepository struct {
	receivedCreateInput Record
	create              func(ctx context.Context, txn Record) (string, error)
	receivedFindInput   string
	findByID            func(ctx context.Context, id string) (*Retrieve, error)
}

func (s *stubRepository) Create(ctx context.Context, txn Record) (string, error) {
	s.receivedCreateInput = txn
	return s.create(ctx, txn)
}

func (s *stubRepository) FindByID(ctx context.Context, id string) (*Retrieve, error) {
	s.receivedFindInput = id
	return s.findByID(ctx, id)
}

func TestService_Create(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"

	mockRepo := &stubRepository{
		create: func(ctx context.Context, txn Record) (string, error) {
			return id, nil
		},
	}

	mockIDGen := func() string {
		return id
	}

	input := Request{
		Description: "food",
		Amount:      20.47,
	}

	want := Record{
		ID:          id,
		Description: "food",
		Amount:      20.47,
	}

	svc := NewService(mockRepo, mockIDGen)
	got, gotErr := svc.Create(context.Background(), input)
	assert.NoError(t, gotErr)
	assert.Equal(t, id, got)
	assert.Equal(t, want, mockRepo.receivedCreateInput)
}

func TestService_Create_Error(t *testing.T) {
	someErr := errors.New("some error")

	testCases := map[string]struct {
		input    Request
		mockRepo *stubRepository
		wantErr  error
	}{
		"validation error": {
			input: Request{
				Description: "food",
				Amount:      -5,
			},
			mockRepo: &stubRepository{
				create: func(ctx context.Context, txn Record) (string, error) {
					return "", nil
				},
			},
			wantErr: httpresponse.ErrValidation,
		},
		"repository error": {
			input: Request{
				Description: "food",
				Amount:      20.47,
			},
			mockRepo: &stubRepository{
				create: func(ctx context.Context, txn Record) (string, error) {
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

			svc := NewService(tc.mockRepo, mockIDGen)
			got, gotErr := svc.Create(context.Background(), tc.input)
			assert.Empty(t, got)
			assert.ErrorContains(t, gotErr, tc.wantErr.Error())
		})
	}
}

func TestService_Get(t *testing.T) {
	id := "b62a64c9-0008-4148-99f6-9c8086a1dd42"

	mockRepo := &stubRepository{
		findByID: func(ctx context.Context, id string) (*Retrieve, error) {
			return &Retrieve{
				ID:          id,
				Description: "food",
				CreatedAt:   time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
				Amount:      20.31,
			}, nil
		},
	}

	mockIDGen := func() string {
		return id
	}

	svc := NewService(mockRepo, mockIDGen)
	got, gotErr := svc.Get(context.Background(), id)
	assert.NoError(t, gotErr)

	want := &Retrieve{
		ID:          id,
		Description: "food",
		CreatedAt:   time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
		Amount:      20.31,
	}
	assert.Equal(t, want, got)
	assert.Equal(t, id, mockRepo.receivedFindInput)
}

func TestService_Get_Error(t *testing.T) {
	someErr := errors.New("some error")

	testCases := map[string]struct {
		input    string
		mockRepo *stubRepository
		wantErr  error
	}{
		"validation error": {
			input: "invalid-uuid",
			mockRepo: &stubRepository{
				findByID: func(ctx context.Context, id string) (*Retrieve, error) {
					return nil, nil
				},
			},
			wantErr: httpresponse.ErrValidation,
		},
		"repository error": {
			input: "b62a64c9-0008-4148-99f6-9c8086a1dd42",
			mockRepo: &stubRepository{
				findByID: func(ctx context.Context, id string) (*Retrieve, error) {
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

			svc := NewService(tc.mockRepo, mockIDGen)
			got, gotErr := svc.Get(context.Background(), tc.input)
			assert.Nil(t, got)
			assert.ErrorContains(t, gotErr, tc.wantErr.Error())
		})
	}
}
