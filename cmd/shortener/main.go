package main

import (
	"errors"
	"fmt"
	"github.com/alexeilarionov/url-shortener/internal/app/config"
	"github.com/alexeilarionov/url-shortener/internal/app/handlers"
	"github.com/alexeilarionov/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	cfg := config.Load()
	var store storage.Storage

	var err error
	switch cfg.Storage {
	case "memory":
		store = storage.NewInMemoryStorage()
	default:
		err = errors.New("unknown storage type: " + cfg.Storage)
	}

	h := &handlers.Handler{
		ShortURLAddr: cfg.ShortURLAddr,
		Store:        store,
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.ShortenerHandler)
		r.Get("/{id}", h.UnshortenerHandler)
	})

	fmt.Println("Running server on", cfg.StartAddr)
	err = http.ListenAndServe(cfg.StartAddr, r)
	if err != nil {
		panic(err)
	}
}
