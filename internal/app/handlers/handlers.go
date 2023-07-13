package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

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

	ShortenRequest struct {
		URL string `json:"url"`
	}

	ShortenResponse struct {
		Result string `json:"result"`
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

	err = h.Store.Store(storage.ShortenedData{
		UUID:        uuid.New().String(),
		ShortURL:    encoded,
		OriginalURL: string(body),
	})
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
	data, err := h.Store.Get(shortenerID)
	if err != nil {
		http.Error(res, "Bad request", http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", data.OriginalURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) JSONShortenerHandler(w http.ResponseWriter, r *http.Request) {
	var sr ShortenRequest
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &sr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoded := hashutil.Encode([]byte(sr.URL))
	err = h.Store.Store(storage.ShortenedData{
		UUID:        uuid.New().String(),
		ShortURL:    encoded,
		OriginalURL: sr.URL,
	})
	if err != nil {
		return
	}

	url := h.ShortURLAddr + "/" + encoded

	resp, err := json.Marshal(ShortenResponse{Result: url})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
