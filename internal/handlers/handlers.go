package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"shortener/config"
	"shortener/internal/auth"
	"shortener/internal/models"
	"shortener/internal/short"
	"shortener/internal/storage"
	"shortener/internal/storage/db"
)

func CreateShortURL(ctx context.Context, w http.ResponseWriter, r *http.Request, cfg config.Config, store storage.Storage, logger *zap.SugaredLogger) {
	rCtx := r.Context()
	userID := rCtx.Value(auth.ContextKey("userID"))
	statusCode := http.StatusCreated
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		logger.Errorw("Couldn't parse url", "error", err)
	}

	if _, err := url.ParseRequestURI(string(body)); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		logger.Errorw("Provided url is not valid", "error", err)
		return
	}

	hash := short.URL(body)
	err = store.Put(ctx, hash, string(body), userID.(string))
	alreadySaved := errors.Is(err, db.ErrorConflict)
	if err != nil && !alreadySaved {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("Couldn't write url to storage", "error", err)
		return
	}

	if alreadySaved {
		statusCode = http.StatusConflict
	}

	shortURL, err := url.JoinPath(cfg.BaseURL, hash)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("Couldn't write url to storage", "error", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("CreateShortURL Handler response", "error", err)
		return
	}
}

func Shorten(ctx context.Context, w http.ResponseWriter, r *http.Request, cfg config.Config, store storage.Storage, logger *zap.SugaredLogger) {
	var req models.Request
	statusCode := http.StatusCreated
	rCtx := r.Context()
	userID := rCtx.Value(auth.ContextKey("userID"))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("can't decode url", "error", err)
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		logger.Errorw("provided url is no valid", "error", err)
		return
	}

	hash := short.URL([]byte(req.URL))
	fmt.Println(hash)
	err := store.Put(ctx, hash, req.URL, userID.(string))
	alreadySaved := errors.Is(err, db.ErrorConflict)
	if err != nil && !alreadySaved {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("can't save url in db", "error", err)
		return
	}

	if alreadySaved {
		statusCode = http.StatusConflict
	}

	shortURL, err := url.JoinPath(cfg.BaseURL, hash)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("can't create url", "error", err)
		return
	}

	resp := models.Response{
		Result: shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("Can't encode url", "error", err)
		return
	}
}

func GetShortURL(w http.ResponseWriter, r *http.Request, id string, store storage.Storage, logger *zap.SugaredLogger) {
	v, err := store.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		logger.Errorw("Can't find shorten url", "error", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func ShortenBatch(ctx context.Context, w http.ResponseWriter, r *http.Request, cfg config.Config, store storage.Storage, logger *zap.SugaredLogger) {
	rCtx := r.Context()
	userID := rCtx.Value(auth.ContextKey("userID"))

	var urls models.BatchRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&urls); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("can't decode urls batch", "error", err)
		return
	}

	if len(urls) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		logger.Errorw("body should contains url")
		return
	}

	var dbBatch []models.URLItem
	for _, u := range urls {
		item := models.URLItem{
			CorrelationID: u.CorrelationID,
			OriginalURL:   u.OriginalURL,
			ShortURL:      short.URL([]byte(u.OriginalURL)),
		}
		dbBatch = append(dbBatch, item)
	}

	err := store.Batch(ctx, dbBatch, userID.(string))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("Can't save urls in storage", "error", err)
		return
	}

	var response models.BatchResponse
	for _, i := range dbBatch {
		shortURL, err := url.JoinPath(cfg.BaseURL, short.URL([]byte(i.OriginalURL)))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			logger.Errorw("Can't create url", "error", err)
			return
		}

		item := models.BatchResponseItem{
			CorrelationID: i.CorrelationID,
			ShortURL:      shortURL,
		}

		response = append(response, item)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("Can't encode url", "error", err)
		return
	}
}

func PingDB(w http.ResponseWriter, r *http.Request, store storage.Storage, logger *zap.SugaredLogger) {
	if err := store.Ping(r.Context()); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("Can't ping db", "error", err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetAllURLs(ctx context.Context, w http.ResponseWriter, r *http.Request, store storage.Storage, logger *zap.SugaredLogger) {
	rCtx := r.Context()
	userID := rCtx.Value(auth.ContextKey("userID"))

	urls, err := store.GetAllURLs(ctx, userID.(string))
	logger.Info(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("failed to get user urls", "err", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		if err := enc.Encode(urls); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			logger.Errorw("error encoding response", "err", err)
			return
		}
	}
}
