package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/alexeilarionov/url-shortener/internal/app/logger"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandler(next http.Handler) http.Handler {
	gzipFn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gw := gzipWriter{ResponseWriter: w}
		gz, err := gzip.NewWriterLevel(&gw, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		gw.Writer = gz
		next.ServeHTTP(&gw, r)

		contentType := gw.Header().Get("Content-Type")
		if contentType != "application/json" && contentType != "text/html" {
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
	}
	return http.HandlerFunc(gzipFn)
}

func UnzipRequest(next http.Handler) http.Handler {
	gzipFn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusInternalServerError)
				return
			}
			defer reader.Close()

			r.Body = http.MaxBytesReader(w, reader, r.ContentLength)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(gzipFn)
}

func RequestLogger(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()

		logger.Log.Info("got incoming HTTP request",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.Int("duration", int(duration)),
		)

	}
	return http.HandlerFunc(logFn)
}

func ResponseLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := NewLoggingResponseWriter(w)
		next.ServeHTTP(srw, r)
		logger.Log.Info("response sent",
			zap.Int("status", srw.Status),
			zap.Int("size", srw.Size),
		)
	})
}
