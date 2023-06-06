package main

import (
	"github.com/alexeilarionov/url-shortener/handlers"
	"github.com/alexeilarionov/url-shortener/internal/app/data"
	"net/http"
)

func main() {
	data.InitData()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Handler)
	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
