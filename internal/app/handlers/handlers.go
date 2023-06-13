package handlers

import (
	"github.com/alexeilarionov/url-shortener/internal/app/hashutil"
	"github.com/alexeilarionov/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type Handler struct {
	ShortUrlAddr string
	Store        storage.Storage
}

func (h *Handler) ShortenerHandler(res http.ResponseWriter, req *http.Request) {
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

	err = h.Store.Store(encoded, string(body))
	if err != nil {
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	url := h.ShortUrlAddr + "/" + encoded
	_, err = res.Write([]byte(url))
	if err != nil {
		http.Error(res, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UnshortenerHandler(res http.ResponseWriter, req *http.Request) {
	shortenerID := chi.URLParam(req, "id")
	url, err := h.Store.Get(shortenerID)
	if err != nil {
		http.Error(res, "Bad request", http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
