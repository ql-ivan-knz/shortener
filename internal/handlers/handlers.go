package handlers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
	"shortener/internal/short"
	"shortener/internal/storage"
)

func CreateShortHandler(w http.ResponseWriter, r *http.Request) {
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
	id := chi.URLParam(r, "id")

	v, err := storage.Get(string(id))
	if err != nil {
		http.Error(w, "Couldn't find"+string(v), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
