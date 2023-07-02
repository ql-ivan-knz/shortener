package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"shortener/config"
	"shortener/internal/handlers"
	"shortener/internal/storage"
)

type Handlers struct {
	config  config.Config
	storage storage.Storage
	logger  *zap.SugaredLogger
	ctx     context.Context
}

func NewHandlers(ctx context.Context, cfg config.Config, storage storage.Storage, l *zap.SugaredLogger) *Handlers {
	return &Handlers{
		config:  cfg,
		storage: storage,
		logger:  l,
		ctx:     ctx,
	}
}

func (h *Handlers) createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	handlers.CreateShortURL(h.ctx, w, r, h.config, h.storage, h.logger)
}

func (h *Handlers) shortenHandler(w http.ResponseWriter, r *http.Request) {
	handlers.Shorten(h.ctx, w, r, h.config, h.storage, h.logger)
}

func (h *Handlers) getShortURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	handlers.GetShortURL(w, r, string(id), h.storage, h.logger)
}

func (h *Handlers) shortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	handlers.ShortenBatch(h.ctx, w, r, h.config, h.storage, h.logger)
}

func (h *Handlers) getAllURLs(w http.ResponseWriter, r *http.Request) {
	handlers.GetAllURLs(h.ctx, w, r, h.storage, h.logger)
}

func (h *Handlers) pingDBHandler(w http.ResponseWriter, r *http.Request) {
	handlers.PingDB(w, r, h.storage, h.logger)
}
