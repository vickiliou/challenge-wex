package config

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/vickiliou/challenge-wex/internal/gateway"
	"github.com/vickiliou/challenge-wex/internal/httphandler"
	"github.com/vickiliou/challenge-wex/internal/repository"
	"github.com/vickiliou/challenge-wex/internal/transaction"
)

// SetupRouter creates and configures the HTTP router for the application.
func SetupRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	gw := gateway.NewGateway(&http.Client{})
	repo := repository.NewRepository(db)
	svc := transaction.NewService(repo, gw, uuid.NewString)
	h := httphandler.NewHandler(svc)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/v1/transactions", h.Store)
	r.Get("/v1/transactions/{id}", h.Retrieve)

	return r
}
