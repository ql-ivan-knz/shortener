package server

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"shortener/config"
	"shortener/internal/handlers"
	"shortener/internal/logger"
	"shortener/internal/storage"
)

type payload struct {
	config  config.Config
	storage *storage.Storage
}

func StartServer() {
	p := payload{storage: storage.NewStorage(), config: config.GetConfig()}
	r := chi.NewRouter()
	logger.Initialize()

	r.Post("/", logger.WithLogging(p.createShortURLHandler()))
	r.Get("/{id}", logger.WithLogging(p.getShortURLHandler()))

	logger.Log.Info("Server started at", zap.String("address", p.config.ServerAddr))
	err := http.ListenAndServe(p.config.ServerAddr, r)
	if err != nil {
		logger.Log.Fatal(err.Error())
		return
	}
}

func (p *payload) createShortURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateShortURL(w, r, p.config, p.storage)
	}
}

func (p *payload) getShortURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handlers.GetShortURL(w, r, string(id), p.storage)
	}
}
