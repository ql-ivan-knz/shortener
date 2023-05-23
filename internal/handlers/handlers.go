package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"shortener/config"
	"shortener/internal/short"
	"shortener/internal/storage"
)

func CreateShortURL(w http.ResponseWriter, r *http.Request, cfg config.Config, store *storage.Storage) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Couldn't parse url", http.StatusBadRequest)
	}

	if _, err := url.ParseRequestURI(string(body)); err != nil {
		http.Error(w, "Provided url is not valid", http.StatusBadRequest)
		return
	}

	hash := short.URL(body)
	store.Set(hash, string(body))
	shortURL, err := url.JoinPath(cfg.BaseURL, hash)
	if err != nil {
		http.Error(w, "Couldn't generate url", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetShortURL(w http.ResponseWriter, r *http.Request, id string, store *storage.Storage) {
	v, ok := store.Get(id)
	if !ok {
		http.Error(w, fmt.Sprintf("Couldn't find %s", v), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
