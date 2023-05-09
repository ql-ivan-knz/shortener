package handlers

import (
	"io"
	"net/http"
	"net/url"
	"shortener/internal/short"
	"shortener/internal/storage"
	"strings"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		CreateShortHandler(w, r)
		return
	}

	GetShortHandler(w, r)
}

func CreateShortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Couldn't parse url", http.StatusBadRequest)
	}

	if _, err := url.ParseRequestURI(string(body)); err != nil {
		http.Error(w, "Provided url is not valid", http.StatusBadRequest)
		return
	}

	hash := short.URL(body)
	storage.Set(hash, string(body))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("http://localhost:8080/" + hash))
}

func GetShortHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	if len(pathParts) > 1 {
		http.Error(w, "expected url be /:id", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusBadRequest)
		return
	}

	key := pathParts[0]

	v, err := storage.Get(string(key))
	if err != nil {
		http.Error(w, "Couldn't find"+string(v), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
