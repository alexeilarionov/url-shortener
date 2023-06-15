package main

import (
	"fmt"
	"github.com/alexeilarionov/url-shortener/internal/app/config"
	"github.com/alexeilarionov/url-shortener/internal/app/handlers"
	"github.com/alexeilarionov/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	cfg := config.Load()

	store := storage.NewStorage(cfg.StorageType)

	h := &handlers.Handler{
		ShortURLAddr: cfg.ShortURLAddr,
		Store:        store,
	}

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.ShortenerHandler)
		r.Get("/{id}", h.UnshortenerHandler)
	})

	fmt.Println("Running server on", cfg.StartAddr)
	err := http.ListenAndServe(cfg.StartAddr, r)
	if err != nil {
		panic(err)
	}
}
