package main

import (
	"crud-users-in-memory/api"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
		return
	}
	slog.Info("all systems offline")
}

func run() error {
	db := newApplication()

	handler := api.NewHandler(db)
	s := http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func newApplication() *api.Application {
	return &api.Application{
		Data: make(map[api.Id]api.User),
	}
}
