package server

import (
	"go.uber.org/zap"
	"net/http"
	"shortener/config"
	"shortener/internal/auth"
	"shortener/internal/middleware/compress"
	"shortener/internal/middleware/logger"
)

type Middleware struct {
	logger *zap.SugaredLogger
	cfg    config.Config
}

func NewMiddleware(lg *zap.SugaredLogger, config config.Config) *Middleware {
	return &Middleware{
		logger: lg,
		cfg:    config,
	}
}

func (m *Middleware) withLogging(h http.Handler) http.Handler {
	return logger.WithLogging(h, m.logger)
}

func (m *Middleware) withCompressing(h http.Handler) http.Handler {
	return compress.Gzip(h, m.logger)
}

func (m *Middleware) withAuth(h http.Handler) http.Handler {
	return auth.WithAuth(h, m.cfg, m.logger)
}
