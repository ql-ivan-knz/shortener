package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"shortener/config"
	"shortener/internal/models"
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

func Shorten(w http.ResponseWriter, r *http.Request, cfg config.Config, store *storage.Storage) {
	var req models.Request

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Couldn't read url", http.StatusInternalServerError)
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, "Provided url is not valid", http.StatusBadRequest)
		return
	}

	hash := short.URL([]byte(req.URL))
	store.Set(hash, req.URL)
	shortURL, err := url.JoinPath(cfg.BaseURL, hash)
	if err != nil {
		http.Error(w, "Couldn't generate url", http.StatusInternalServerError)
		return
	}

	resp := models.Response{
		Result: shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Couldn't send shorten url", http.StatusInternalServerError)
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
