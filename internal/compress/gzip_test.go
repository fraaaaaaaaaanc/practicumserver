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
	strg := storage.NewStorage()
	log, _ := logger.NewZapLogger(false)
	hndlrs := handlers.NewHandlers(strg, log.Logger, "http://localhost:8080",
		"host=localhost user=postgres password=1234 dbname=video sslmode=disable",
		"/tmp/short-url-db.json")
	//"C:\\Users\\frant\\go\\go1.21.0\\bin\\pkg\\mod\\github.com\\fraaaaaaaaaanc\\practicumserver\\internal\\tmp\\short-url-db.json")

	adapter := HandlerFuncAdapter(hndlrs.PostRequestAPIShorten)
	newHandler := MiddlewareGzipHandleFunc(log.Logger)
	handler := newHandler(adapter)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	requestBody := `{"url": "http://test"}`

	successBody := `{"result": "http://localhost:8080/test"}`

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
		assert.JSONEq(t, successBody, string(b))
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
		assert.JSONEq(t, successBody, string(b))
	})
}
