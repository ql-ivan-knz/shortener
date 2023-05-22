package handlers

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateShortURL(t *testing.T) {
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
			body:                strings.NewReader("https://github.com"),
			expectedCode:        http.StatusCreated,
			expectedContentType: "text/plain",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "/", test.body)
			w := httptest.NewRecorder()
			CreateShortURL(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.expectedCode)

			if test.expectedContentType != "" {
				assert.Equal(t, strings.Split(res.Header.Get("content-type"), ";")[0], test.expectedContentType)
			}
		})
	}
}

func TestGetShortURL(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		{
			name:         "returns 400 status code when url has additional params",
			method:       http.MethodGet,
			path:         "/id/wrong/url",
			expectedCode: http.StatusBadRequest,
		}, {
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
			GetShortURL(w, r)

			assert.Equal(t, test.expectedCode, w.Code)
		})
	}
}
