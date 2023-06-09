package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/config"
	"shortener/internal/middleware/logger"
	"shortener/internal/models"
	"shortener/internal/storage"
	"strings"
	"testing"
)

func TestCreateShortURL(t *testing.T) {
	cfg := config.Config{
		ServerAddr:      "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
	}
	store, _ := storage.NewStorage(cfg)

	tests := []struct {
		name                string
		method              string
		body                io.Reader
		expectedCode        int
		expectedContentType string
		expectedBody        string
	}{
		{
			name:                "returns 400 status code if method not POST",
			method:              http.MethodGet,
			expectedCode:        http.StatusBadRequest,
			expectedContentType: "",
		},
		{
			name:                "returns 400 status code if cannot parse body",
			method:              http.MethodPost,
			body:                strings.NewReader("z234dkj;ak"),
			expectedCode:        http.StatusBadRequest,
			expectedContentType: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "/", test.body)
			w := httptest.NewRecorder()
			l, _ := logger.NewLogger()
			CreateShortURL(context.Background(), w, r, cfg, store, l)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.expectedCode)

			if test.expectedContentType != "" {
				assert.Equal(t, strings.Split(res.Header.Get("content-type"), ";")[0], test.expectedContentType)
			}
		})
	}
}

func TestCreateShortURLJSON(t *testing.T) {
	cfg := config.Config{
		ServerAddr:      "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
	}
	store, _ := storage.NewStorage(cfg)

	tests := []struct {
		name         string
		method       string
		body         models.Request
		expectedCode int
		expectedBody string
	}{
		{
			name:         "returns error code if provided bad url",
			method:       http.MethodPost,
			body:         models.Request{URL: ""},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			r := httptest.NewRequest(test.method, "/shorten", bytes.NewReader(body))
			w := httptest.NewRecorder()
			l, _ := logger.NewLogger()
			Shorten(context.Background(), w, r, cfg, store, l)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.expectedCode)
		})
	}
}
