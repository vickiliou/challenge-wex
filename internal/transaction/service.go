package transaction

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/vickiliou/challenge-wex/gateway"
	"github.com/vickiliou/challenge-wex/internal/httpresponse"
)

type repository interface {
	Create(ctx context.Context, txn Transaction) (string, error)
	FindByID(ctx context.Context, id string) (*Transaction, error)
}

type gatewayExchangeRate interface {
	GetExchangeRate(txnDate time.Time) (*gateway.CurrencyExchangeRate, error)
}

type uuidGenerator func() string

// Service represents the transaction service that encapsulates the business logic related to transactions.
type Service struct {
	repo        repository
	gw          gatewayExchangeRate
	idGenerator uuidGenerator
}

// NewService creates a new instance of the transaction service.
func NewService(repo repository, gw gatewayExchangeRate, idGenerator uuidGenerator) *Service {
	return &Service{
		repo:        repo,
		gw:          gw,
		idGenerator: idGenerator,
	}
}

// Create creates a new transaction based on user input.
func (s *Service) Create(ctx context.Context, input RecordRequest) (string, error) {
	if err := input.validate(); err != nil {
		return "", fmt.Errorf("%w: %s", httpresponse.ErrValidation, err.Error())
	}

	txn := Transaction{
		ID:              s.idGenerator(),
		Description:     input.Description,
		TransactionDate: input.TransactionDate,
		Amount:          input.Amount,
	}

	return s.repo.Create(ctx, txn)
}

// Get retrieves a transaction by its ID.
func (s *Service) Get(ctx context.Context, id string) (*RetrieveResponse, error) {
	if ok := isValidUUID(id); !ok {
		return nil, fmt.Errorf("%w: id must be a valid UUID", httpresponse.ErrValidation)
	}

	txn, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error calling database: %s", err.Error())
	}

	exchangeRate, err := s.gw.GetExchangeRate(txn.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("error calling gateway: %s", err.Error())
	}

	exchangeRateFloat, err := strconv.ParseFloat(exchangeRate.ExchangeRate, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting exchange rate: %s", err.Error())
	}

	return &RetrieveResponse{
		ID:              txn.ID,
		Description:     txn.Description,
		TransactionDate: txn.TransactionDate,
		OriginalAmount:  txn.Amount,
		ExchangeRate:    exchangeRateFloat,
		ConvertedAmount: roundTwoDecimal(exchangeRateFloat * txn.Amount),
	}, nil
}
