package handlers

import (
	"github.com/alexeilarionov/url-shortener/internal/app/data"
	"github.com/alexeilarionov/url-shortener/internal/app/hashutil"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func ShortenerHandler(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(res, "Empty request body", http.StatusBadRequest)
		return
	}
	encoded := hashutil.Encode(body)

	err = data.AddOrUpdateData(encoded, string(body))
	if err != nil {
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	url := "http://localhost:8080/" + encoded
	_, err = res.Write([]byte(url))
	if err != nil {
		http.Error(res, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func UnshortenerHandler(res http.ResponseWriter, req *http.Request) {
	shortenerID := chi.URLParam(req, "id")
	url, err := data.GetDataByKey(shortenerID)
	if err != nil {
		http.Error(res, "Bad request", http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
