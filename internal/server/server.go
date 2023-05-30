package server

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"shortener/config"
	"shortener/internal/handlers"
	"shortener/internal/middleware/compress"
	"shortener/internal/middleware/logger"
	"shortener/internal/storage"
)

type payload struct {
	config  config.Config
	storage storage.Storage
}

func StartServer() {
	cfg := config.GetConfig()
	p := payload{storage: storage.NewStorage(cfg.FileStoragePath), config: cfg}
	r := chi.NewRouter()
	logger.Initialize()

	r.Post("/", logger.WithLogging(compress.Gzip(p.createShortURLHandler)))
	r.Post("/api/shorten", logger.WithLogging(compress.Gzip(p.shortenHandler)))
	r.Get("/{id}", logger.WithLogging(compress.Gzip(p.getShortURLHandler)))

	logger.Log.Info("Server started at", zap.String("address", p.config.ServerAddr), zap.String("storage path", cfg.FileStoragePath))
	err := http.ListenAndServe(p.config.ServerAddr, r)
	if err != nil {
		logger.Log.Fatal(err.Error())
		return
	}
}

func (p *payload) createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	handlers.CreateShortURL(w, r, p.config, p.storage)
}

func (p *payload) shortenHandler(w http.ResponseWriter, r *http.Request) {
	handlers.Shorten(w, r, p.config, p.storage)
}

func (p *payload) getShortURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	handlers.GetShortURL(w, r, string(id), p.storage)
}
