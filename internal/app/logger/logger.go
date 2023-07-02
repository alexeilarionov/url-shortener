package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/alexeilarionov/url-shortener/internal/app/handlers"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

func RequestLogger(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()

		Log.Info("got incoming HTTP request",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.Int("duration", int(duration)),
		)

	}
	return http.HandlerFunc(logFn)
}

func ResponseLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := handlers.NewLoggingResponseWriter(w)
		next.ServeHTTP(srw, r)
		Log.Info("response sent",
			zap.Int("status", srw.Status),
			zap.Int("size", srw.Size),
		)
	})
}
