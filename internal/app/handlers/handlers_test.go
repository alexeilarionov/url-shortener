package handlers

import (
	"context"
	"github.com/alexeilarionov/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenerHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "positive shortener test",
			body: "test.com",
			want: want{
				code:        201,
				response:    `http://localhost:8080/yXwbNnH`,
				contentType: "text/plain",
			},
		},
		{
			name: "shortener empty body test",
			body: "",
			want: want{
				code:        400,
				response:    "Empty request body\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	h := &Handler{
		ShortUrlAddr: "http://localhost:8080",
		Store:        storage.NewInMemoryStorage(),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.ShortenerHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestUnshortenerHandler(t *testing.T) {
	type want struct {
		code     int
		location string
		body     string
	}
	tests := []struct {
		name   string
		target string
		want   want
	}{
		{
			name:   "positive unshortener test",
			target: "yXwbNnH",
			want: want{
				code:     307,
				location: "test.com",
				body:     "",
			},
		},
		{
			name:   "unshortener empty id",
			target: "",
			want: want{
				code:     400,
				location: "",
				body:     "Bad request\n",
			},
		},
		{
			name:   "unshortener not exist id",
			target: "123456",
			want: want{
				code:     400,
				location: "",
				body:     "Bad request\n",
			},
		},
	}
	h := &Handler{
		ShortUrlAddr: "http://localhost:8080",
		Store:        storage.NewInMemoryStorage(),
	}
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("test.com"))
	// создаём новый Recorder
	w := httptest.NewRecorder()
	h.ShortenerHandler(w, request)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/"+test.target, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.target)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			h.UnshortenerHandler(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.body, string(resBody))
			assert.Equal(t, test.want.location, res.Header.Get("Location"))
		})
	}
}
