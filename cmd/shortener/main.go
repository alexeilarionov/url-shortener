package main

import (
	"github.com/alexeilarionov/url-shortener/handlers"
	"github.com/alexeilarionov/url-shortener/internal/app/data"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	data.InitData()

	r := chi.NewRouter()
	r.Get("/{id}", handlers.UnshortenerHandler)
	r.Post("/", handlers.ShortenerHandler)

	err := http.ListenAndServe(`localhost:8080`, r)
	if err != nil {
		panic(err)
	}
}
