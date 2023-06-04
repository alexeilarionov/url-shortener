package main

import (
	"github.com/alexeilarionov/url-shortener/internal/app"
	"io"
	"net/http"
)

var (
	data = make(map[string]string)
)

func handler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		if req.Header.Get("Content-Type") != "text/plain" {
			http.Error(res, "Invalid content type", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, "Failed to read request body", http.StatusBadRequest)
			return
		}
		if len(body) == 0 {
			http.Error(res, "Empty request body", http.StatusBadRequest)
			return
		}
		encoded := app.Encode(body)
		data[encoded] = string(body)
		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(encoded))

	} else if req.Method == http.MethodGet {
		path := req.URL.Path
		if len(path) < 2 {
			http.Error(res, "Bad request", http.StatusBadRequest)
			return
		}
		path = path[1:]
		if _, ok := data[path]; ok {
			http.Redirect(res, req, data[path], http.StatusTemporaryRedirect)
		}
	} else {
		http.Error(res, "Not supported", http.StatusBadRequest)
		return
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
