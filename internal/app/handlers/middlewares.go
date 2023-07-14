package handlers

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/alexeilarionov/url-shortener/internal/app/logger"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	Status int
	Size   int
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gw            *gzip.Writer
	buffer        bytes.Buffer
	statusCode    int
	headerWritten bool
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.buffer.Write(b)
}

func (w *gzipResponseWriter) flushHeader() {
	if !w.headerWritten {
		contentType := w.Header().Get("Content-Type")
		if contentType == "application/json" || contentType == "text/html" {
			w.Header().Set("Content-Encoding", "gzip")
			w.ResponseWriter.WriteHeader(w.statusCode)
			w.gw = gzip.NewWriter(w.ResponseWriter)
		} else {
			w.ResponseWriter.WriteHeader(w.statusCode)
		}
		w.headerWritten = true
	}
}

func (w *gzipResponseWriter) flushBody() {
	if w.gw != nil {
		w.gw.Write(w.buffer.Bytes())
		w.gw.Close()
	} else {
		w.ResponseWriter.Write(w.buffer.Bytes())
	}
}

func GzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		grw := &gzipResponseWriter{ResponseWriter: w, statusCode: 200}
		next.ServeHTTP(grw, r)

		grw.flushHeader()
		grw.flushBody()
	})
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
