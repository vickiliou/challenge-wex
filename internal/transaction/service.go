package transaction

import (
	"context"
	"fmt"

	"github.com/vickiliou/challenge-wex/internal/httpresponse"
)

type repository interface {
	Create(ctx context.Context, txn Record) (string, error)
	FindByID(ctx context.Context, id string) (*Retrieve, error)
}

type uuidGenerator func() string

// Service represents the transaction service that encapsulates the business logic related to transactions.
type Service struct {
	repo        repository
	idGenerator uuidGenerator
}

// NewService creates a new instace of the transaction Service.
func NewService(repo repository, idGenerator uuidGenerator) *Service {
	return &Service{
		repo:        repo,
		idGenerator: idGenerator,
	}
}

// Create creates a new transaction based on user input.
func (s *Service) Create(ctx context.Context, input Request) (string, error) {
	if err := input.validate(); err != nil {
		return "", fmt.Errorf("%w: %s", httpresponse.ErrValidation, err.Error())
	}

	txn := Record{
		ID:          s.idGenerator(),
		Description: input.Description,
		Amount:      input.Amount,
	}

	return s.repo.Create(ctx, txn)
}

// Get retrieves a transaction by its ID.
func (s *Service) Get(ctx context.Context, id string) (*Retrieve, error) {
	if ok := isValidUUID(id); !ok {
		return nil, fmt.Errorf("%w: id must be a valid UUID", httpresponse.ErrValidation)
	}

	return s.repo.FindByID(ctx, id)
}
