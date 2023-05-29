package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/config"
	"shortener/internal/handlers"
	"shortener/internal/storage"
	"testing"
)

func TestGzip(t *testing.T) {
	var storeMock storage.Storage = make(map[string]string)
	configMock := config.Config{
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost:8080",
	}
	requestBody := `{ "url": "https://github.com" }`
	successBody := `{ "result": "http://localhost:8080/3097fca9" }`

	handler := http.HandlerFunc(Gzip(func(w http.ResponseWriter, r *http.Request) {
		handlers.Shorten(w, r, configMock, &storeMock)
	}))

	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		assert.NoError(t, err)
		err = zb.Close()
		assert.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, successBody, string(b))

	})

	t.Run("accept_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		assert.NoError(t, err)

		b, err := io.ReadAll(zr)
		assert.NoError(t, err)

		assert.JSONEq(t, successBody, string(b))
	})
}