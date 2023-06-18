package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"shortener/config"
	"shortener/internal/models"
	"shortener/internal/short"
	"shortener/internal/storage"
	"shortener/internal/storage/db"
)

func CreateShortURL(w http.ResponseWriter, r *http.Request, cfg config.Config, store storage.Storage) {
	statusCode := http.StatusCreated
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Couldn't parse url", http.StatusBadRequest)
	}

	if _, err := url.ParseRequestURI(string(body)); err != nil {
		http.Error(w, "Provided url is not valid", http.StatusBadRequest)
		return
	}

	hash := short.URL(body)
	err = store.Put(r.Context(), hash, string(body))
	alreadySaved := errors.Is(err, db.ErrorConflict)
	if err != nil && !alreadySaved {
		http.Error(w, "Couldn't write url to storage", http.StatusInternalServerError)
	}

	if alreadySaved {
		statusCode = http.StatusConflict
	}

	shortURL, err := url.JoinPath(cfg.BaseURL, hash)
	if err != nil {
		http.Error(w, "Couldn't generate url", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func Shorten(w http.ResponseWriter, r *http.Request, cfg config.Config, store storage.Storage) {
	var req models.Request
	statusCode := http.StatusCreated

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
	fmt.Println(hash)
	err := store.Put(r.Context(), hash, req.URL)
	alreadySaved := errors.Is(err, db.ErrorConflict)
	if err != nil && !alreadySaved {
		http.Error(w, "Couldn't write url to storage", http.StatusInternalServerError)
	}

	if alreadySaved {
		statusCode = http.StatusConflict
	}
	shortURL, err := url.JoinPath(cfg.BaseURL, hash)
	fmt.Println(shortURL)
	if err != nil {
		http.Error(w, "Couldn't generate url", http.StatusInternalServerError)
		return
	}

	resp := models.Response{
		Result: shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Couldn't send shorten url", http.StatusInternalServerError)
		return
	}
}

func GetShortURL(w http.ResponseWriter, r *http.Request, id string, store storage.Storage) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusBadRequest)
		return
	}

	v, err := store.Get(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't find %s", v), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func ShortenBatch(w http.ResponseWriter, r *http.Request, cfg config.Config, store storage.Storage) {
	var urls models.BatchRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&urls); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err := store.Batch(r.Context(), urls)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var response models.BatchResponse
	for _, u := range urls {
		shortURL, err := url.JoinPath(cfg.BaseURL, short.URL([]byte(u.OriginalURL)))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		response = append(response, models.BatchResponseItem{CorrelationID: u.CorrelationID, ShortURL: shortURL})
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func PingDB(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusBadRequest)
		return
	}

	if err := store.Ping(r.Context()); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
