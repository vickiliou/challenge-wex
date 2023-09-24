package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vickiliou/challenge-wex/internal/httpresponse"
	"github.com/vickiliou/challenge-wex/internal/transaction"
	"golang.org/x/exp/slog"
)

type service interface {
	Create(ctx context.Context, input transaction.RecordRequest) (string, error)
	Get(ctx context.Context, id string) (*transaction.Retrieve, error)
}

// Handler is responsible for handling HTTP requests related to transactions.
type Handler struct {
	svc service
}

// NewHandler creates a new transaction handler with the given service.
func NewHandler(svc service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// Store handles the creation of a new transaction.
func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	var input transaction.RecordRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpresponse.RespondWithError(w, http.StatusBadRequest, err)
		httpresponse.LogError("Error decoding request body", http.StatusBadRequest, err)
		return
	}

	id, err := h.svc.Create(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, httpresponse.ErrValidation):
			httpresponse.RespondWithError(w, http.StatusBadRequest, err)
			httpresponse.LogError("Validation error", http.StatusBadRequest, err)
			return
		default:
			httpresponse.RespondWithError(w, http.StatusInternalServerError, err)
			httpresponse.LogError("Unexpected error", http.StatusInternalServerError, err)
			return
		}
	}

	res := transaction.RecordResponse{
		ID: id,
	}

	httpresponse.RespondJSON(w, http.StatusCreated, res)
	slog.Info("Transaction created successfully", "ID", id)
}

// Retrieve retrieves a transaction by its ID.
func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.svc.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, httpresponse.ErrValidation):
			httpresponse.RespondWithError(w, http.StatusBadRequest, err)
			httpresponse.LogError("Validation error", http.StatusBadRequest, err)
			return
		case errors.Is(err, httpresponse.ErrNotFound):
			httpresponse.RespondWithError(w, http.StatusNotFound, err)
			httpresponse.LogError("Not found", http.StatusNotFound, err)
			return
		default:
			httpresponse.RespondWithError(w, http.StatusInternalServerError, err)
			httpresponse.LogError("Unexpected error", http.StatusInternalServerError, err)
			return
		}
	}

	httpresponse.RespondJSON(w, http.StatusOK, res)
	slog.Info("Transaction retrieved successfully")
}
