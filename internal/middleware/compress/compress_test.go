package compress

import (
	"context"
	"net/http"
	"net/http/httptest"
	"shortener/config"
	"shortener/internal/handlers"
	"shortener/internal/middleware/logger"
	"shortener/internal/storage"
	"testing"
)

func TestGzip(t *testing.T) {
	configMock := config.Config{
		ServerAddr:      "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
	}
	storeMock, _ := storage.NewStorage(configMock)
	l, _ := logger.NewLogger()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.Shorten(context.Background(), w, r, configMock, storeMock, l)
	})

	handler := Gzip(h, l)

	srv := httptest.NewServer(handler)
	defer srv.Close()
}
