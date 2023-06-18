package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"shortener/config"
	"shortener/internal/handlers"
	"shortener/internal/storage"
)

type Handlers struct {
	config  config.Config
	storage storage.Storage
	ctx     context.Context
}

func NewHandlers(cfg config.Config, storage storage.Storage, ctx context.Context) *Handlers {
	return &Handlers{
		config:  cfg,
		storage: storage,
		ctx:     ctx,
	}
}

func (h *Handlers) createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	handlers.CreateShortURL(w, r, h.config, h.storage)
}

func (h *Handlers) shortenHandler(w http.ResponseWriter, r *http.Request) {
	handlers.Shorten(w, r, h.config, h.storage)
}

func (h *Handlers) getShortURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	handlers.GetShortURL(w, r, string(id), h.storage)
}

func (h *Handlers) shortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	handlers.ShortenBatch(w, r, h.config, h.storage)
}

func (h *Handlers) pingDBHandler(w http.ResponseWriter, r *http.Request) {
	handlers.PingDB(w, r, h.storage)
}
