package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"practicumserver/internal/handlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/storage"
	"testing"
)

type HandlerFuncAdapter func(http.ResponseWriter, *http.Request)

func (h HandlerFuncAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//flagURL := "http://localhost:8080"
	//flagPath := "/tmp/short-url-db.json"
	//hndlr := handlers.NewHandlers()
	h(w, r)
}

func TestMiddlewareGzipHandleFunc(t *testing.T) {
	log, _ := logger.NewZapLogger(false)
	strg, _ := storage.NewStorage(log.Logger, "", "")
	hndlrs := handlers.NewHandlers(strg, log.Logger, "http://localhost:8080")

	adapter := HandlerFuncAdapter(hndlrs.PostRequestAPIShorten)
	newHandler := MiddlewareGzipHandleFunc(log.Logger)
	handler := newHandler(adapter)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	requestBody := `{"url": "http://test.com"}`

	t.Run("send_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		assert.NoError(t, err)
		err = zb.Close()
		assert.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL+"/api/shorten", buf)
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		r.Header.Set("Content-Encoding", "gzip")
		r.RequestURI = ""
		resp, err := http.DefaultClient.Do(r)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.NotNil(t, string(b))
	})

	t.Run("accept_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", srv.URL+"/api/shorten", buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := http.DefaultClient.Do(r)
		assert.NoError(t, err)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		assert.NoError(t, err)

		b, err := io.ReadAll(zr)
		assert.NoError(t, err)
		assert.NotNil(t, string(b))
	})
}
