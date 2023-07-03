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
	userID := rCtx.Value(auth.UserIDContextKey)
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
	userID := rCtx.Value(auth.UserIDContextKey)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("can't decode url", "error", err)
		return
	}

	if req.URL == "" {
		http.Error(w, "url should be provided", http.StatusBadRequest)
		return
	}

	hash := short.URL([]byte(req.URL))

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

func GetShortURL(ctx context.Context, w http.ResponseWriter, r *http.Request, id string, store storage.Storage, logger *zap.SugaredLogger) {
	link, err := store.Get(ctx, id)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		logger.Errorw("Can't find shorten url", "error", err)
		return
	}

	if link.OriginalURL == "" {
		http.Error(w, "Link not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("location", link.OriginalURL)
	if link.IsDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("location", link.OriginalURL)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", link.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func ShortenBatch(ctx context.Context, w http.ResponseWriter, r *http.Request, cfg config.Config, store storage.Storage, logger *zap.SugaredLogger) {
	rCtx := r.Context()
	userID := rCtx.Value(auth.UserIDContextKey)

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

func GetAllURLs(ctx context.Context, w http.ResponseWriter, r *http.Request, cfg config.Config, store storage.Storage, logger *zap.SugaredLogger) {
	rCtx := r.Context()
	userID := rCtx.Value(auth.UserIDContextKey)

	urls, err := store.GetAllURLs(ctx, userID.(string))

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("failed to get user urls", "err", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		var response []models.URLItem
		for _, u := range urls {
			shortURL, err := url.JoinPath(cfg.BaseURL, u.ShortURL)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				logger.Errorw("Can't create url", "error", err)
				return
			}

			item := models.URLItem{
				ShortURL:    shortURL,
				OriginalURL: u.OriginalURL,
			}

			response = append(response, item)
		}
		logger.Infow("user id", "userId", userID)
		logger.Infow("urls in", "urls", urls)
		logger.Infow("response", "res", response)

		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		if err := enc.Encode(response); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			logger.Errorw("error encoding response", "err", err)
			return
		}
	}
}

func DeleteURLs(ctx context.Context, writer http.ResponseWriter, request *http.Request, cfg config.Config, str storage.Storage, logger *zap.SugaredLogger) {
	requestContext := request.Context()
	userID := requestContext.Value(auth.UserIDContextKey)
	var req models.DeleteURLsRequest

	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		logger.Errorw("cannot decode request JSON body", "err", err)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
	go func() {
		err := deletingUserUrls(ctx, str, req, userID.(string))
		if err != nil {
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			logger.Errorw("failed to delete URLs", "err", err)
			return
		}
	}()
}

func deletingUserUrls(ctx context.Context, store storage.Storage, urls []string, userID string) error {
	fmt.Println(urls, userID)
	select {
	case <-ctx.Done():
		return fmt.Errorf("deleting was interrupted by context for userID=%s", userID)
	default:
		if err := store.DeleteURLs(ctx, urls, userID); err != nil {
			return fmt.Errorf("failed to delete user urls with userID=%s: %w", userID, err)
		}
		return nil
	}
}
