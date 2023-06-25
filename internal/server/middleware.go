package server

import (
	"go.uber.org/zap"
	"net/http"
	"shortener/internal/middleware/compress"
	"shortener/internal/middleware/logger"
)

type Middleware struct {
	logger *zap.SugaredLogger
}

func NewMiddleware(lg *zap.SugaredLogger) *Middleware {
	return &Middleware{
		logger: lg,
	}
}

func (m *Middleware) withLogging(h http.Handler) http.Handler {
	return logger.WithLogging(h, m.logger)
}

func (m *Middleware) withCompressing(h http.Handler) http.Handler {
	return compress.Gzip(h, m.logger)
}
