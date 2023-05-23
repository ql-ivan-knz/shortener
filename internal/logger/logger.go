package logger

import (
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

var Log *zap.Logger

func Initialize() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	defer logger.Sync()

	Log = logger
}

func WithLogging(h http.Handler) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.Info("REQUEST",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Duration("duration", duration),
		)

		Log.Info("RESPONSE",
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	}

	return http.HandlerFunc(logFn)
}
