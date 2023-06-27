package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(h *Handlers, m *Middleware) error {
	router := chi.NewRouter()

	router.Use(m.withLogging)
	router.Use(m.withCompressing)

	router.Post("/", h.createShortURLHandler)
	router.Post("/api/shorten", h.shortenHandler)
	router.Post("/api/shorten/batch", h.shortenBatchHandler)
	router.Get("/{id}", h.getShortURLHandler)
	router.Get("/ping", h.pingDBHandler)

	m.logger.Infow("Server started at", "address", h.config.ServerAddr)

	err := http.ListenAndServe(h.config.ServerAddr, router)
	if err != nil {
		return err
	}

	return nil
}
