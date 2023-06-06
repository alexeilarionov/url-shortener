package handlers

import (
	"github.com/alexeilarionov/url-shortener/internal/app/data"
	"github.com/alexeilarionov/url-shortener/internal/app/hashutil"
	"io"
	"net/http"
)

func shortenerHandler(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Failed to read request body", http.StatusBadRequest)
		return
	}
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
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
	url := scheme + "://" + req.Host + "/" + encoded
	_, err = res.Write([]byte(url))
	if err != nil {
		http.Error(res, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func unshortenerHandler(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if len(path) < 2 {
		http.Error(res, "Bad request", http.StatusBadRequest)
		return
	}
	path = path[1:]
	url, err := data.GetDataByKey(path)
	if err != nil {
		http.Error(res, "Bad request", http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func RootHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		shortenerHandler(res, req)
	} else if req.Method == http.MethodGet {
		unshortenerHandler(res, req)
	} else {
		http.Error(res, "Not supported", http.StatusBadRequest)
		return
	}
}
