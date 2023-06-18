package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"shortener/internal/middleware/compress"
	"shortener/internal/middleware/logger"
)

func StartServer(h *Handlers) {
	logger.Initialize()
	router := chi.NewRouter()

	router.Post("/", logger.WithLogging(compress.Gzip(h.createShortURLHandler)))
	router.Post("/api/shorten", logger.WithLogging(compress.Gzip(h.shortenHandler)))
	router.Post("/api/shorten/batch", h.shortenBatch)
	router.Get("/{id}", logger.WithLogging(compress.Gzip(h.getShortURLHandler)))
	router.Get("/ping", logger.WithLogging(h.pingDB))

	fmt.Print("started")
	logger.Log.Info("Server started at", zap.String("address", h.config.ServerAddr))
	err := http.ListenAndServe(h.config.ServerAddr, router)
	if err != nil {
		logger.Log.Fatal(err.Error())
		return
	}
}
