package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/alexeilarionov/url-shortener/internal/app/config"
	"github.com/alexeilarionov/url-shortener/internal/app/handlers"
	"github.com/alexeilarionov/url-shortener/internal/app/logger"
	"github.com/alexeilarionov/url-shortener/internal/app/storage"
)

func main() {
	cfg := config.Load()

	store := storage.NewStorage(cfg.StorageType)

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	h := &handlers.Handler{
		ShortURLAddr: cfg.ShortURLAddr,
		Store:        store,
	}

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(logger.ResponseLogger)
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.ShortenerHandler)
		r.Get("/{id}", h.UnshortenerHandler)
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", h.JsonShortenerHandler)
	})

	logger.Log.Info("Running server", zap.String("address", cfg.StartAddr))
	err := http.ListenAndServe(cfg.StartAddr, r)
	if err != nil {
		panic(err)
	}
}
