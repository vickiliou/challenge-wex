package main

import (
	"net/http"
	"os"

	"github.com/vickiliou/challenge-wex/config"
	"github.com/vickiliou/challenge-wex/database"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := database.Setup()
	if err != nil {
		slog.Warn("Failed to open SQLite database", "error", err.Error())
		return
	}
	defer db.Close()

	r := config.SetupRouter(db)

	errCh := make(chan error, 1)

	go func() {
		slog.Info("Starting server")
		if err := http.ListenAndServe(":8082", r); err != nil {
			errCh <- err
		}
	}()

	if err := <-errCh; err != nil {
		slog.Error("Server error", slog.String("error", err.Error()))
	}

}
