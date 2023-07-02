package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/alexeilarionov/url-shortener/internal/app/hashutil"
	"github.com/alexeilarionov/url-shortener/internal/app/storage"
)

type (
	Handler struct {
		ShortURLAddr string
		Store        storage.Storage
	}

	LoggingResponseWriter struct {
		http.ResponseWriter
		Status int
		Size   int
	}
)

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{ResponseWriter: w}
}

func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.Size += size
	return size, err
}

func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.Status = statusCode
}

func (h *Handler) ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}
	encoded := hashutil.Encode(body)

	err = h.Store.Store(encoded, string(body))
	if err != nil {
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	url := h.ShortURLAddr + "/" + encoded
	_, err = w.Write([]byte(url))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
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
