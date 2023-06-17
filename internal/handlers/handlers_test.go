package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/config"
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

	url := "http://github.com"

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
		{
			name:                "returns 201 status code",
			method:              http.MethodPost,
			body:                strings.NewReader(url),
			expectedCode:        http.StatusCreated,
			expectedContentType: "text/plain",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "/", test.body)
			w := httptest.NewRecorder()
			CreateShortURL(w, r, cfg, store)

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
			name:         "returns error code if method not POST",
			method:       http.MethodGet,
			body:         models.Request{URL: "https://github.com"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "returns error code if provided bad url",
			method:       http.MethodPost,
			body:         models.Request{URL: ""},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "returns 200 status code",
			method:       http.MethodPost,
			body:         models.Request{URL: "https://github.com"},
			expectedCode: http.StatusCreated,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			r := httptest.NewRequest(test.method, "/shorten", bytes.NewReader(body))
			w := httptest.NewRecorder()
			Shorten(w, r, cfg, store)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.expectedCode)
		})
	}
}

func TestGetShortURL(t *testing.T) {
	cfg := config.Config{
		ServerAddr:      "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
	}
	store, _ := storage.NewStorage(cfg)
	pathID := "3097fca9"

	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		{
			name:         "returns 400 status code when method not GET",
			method:       http.MethodPost,
			path:         "/id",
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			GetShortURL(w, r, pathID, store)

			assert.Equal(t, test.expectedCode, w.Code)
		})
	}
}
