package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/alexeilarionov/url-shortener/internal/app/config"
	"github.com/alexeilarionov/url-shortener/internal/app/handlers"
	"github.com/alexeilarionov/url-shortener/internal/app/logger"
	"github.com/alexeilarionov/url-shortener/internal/app/storage"
)

func main() {
	cfg := config.Load()

	store := storage.NewStorage(*cfg)

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	h := &handlers.Handler{
		ShortURLAddr: cfg.ShortURLAddr,
		Store:        store,
	}

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(handlers.UnzipRequest)
	r.Use(handlers.GzipHandler)
	r.Use(logger.ResponseLogger)

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.ShortenerHandler)
		r.Get("/{id}", h.UnshortenerHandler)
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", h.JSONShortenerHandler)
	})

	logger.Log.Info("Running server", zap.String("address", cfg.StartAddr))

	// Create a channel to receive OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.StartAddr,
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-stop
	if cfg.StorageType == "file" {
		err := store.(*storage.FileStorage).Save()
		if err != nil {
			panic(err)
		}
	}
	logger.Log.Info("Server stopped")

	err := srv.Shutdown(context.TODO())
	if err != nil {
		panic(err)
	}
}
