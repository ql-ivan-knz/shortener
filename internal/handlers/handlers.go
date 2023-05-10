package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
	"shortener/config"
	"shortener/internal/short"
	"shortener/internal/storage"
)

func CreateShortURL(w http.ResponseWriter, r *http.Request) {
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
	_, _ = w.Write([]byte(fmt.Sprintf("%s/%s", config.BaseURL, hash)))
}

func GetShortURL(w http.ResponseWriter, r *http.Request) {
	id := getIDParam(r)

	v, err := storage.Get(string(id))
	if err != nil {
		http.Error(w, "Couldn't find"+string(v), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func getIDParam(r *http.Request) string {
	return chi.URLParam(r, "id")
}
