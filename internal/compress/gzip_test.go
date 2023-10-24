package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"practicumserver/internal/cookie"
	handlers "practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	storage "practicumserver/internal/storage/pg"
	"testing"
)

type HandlerFuncAdapter func(http.ResponseWriter, *http.Request)

func (h HandlerFuncAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

func createHTTPAuthClient(log *zap.Logger, srv *httptest.Server) (*http.Client, error) {
	jar, _ := cookiejar.New(nil)
	token, err := cookie.BuildJWTString("testUserID")
	if err != nil {
		log.Error("Error:", zap.Error(err))
		return nil, err
	}
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		MaxAge:   7200,
		HttpOnly: true,
	}
	c := make([]*http.Cookie, 1)
	c[0] = cookie
	urlStr := srv.URL
	parse, err := url.Parse(urlStr)
	if err != nil {
		log.Error("Error:", zap.Error(err))
		return nil, err
	}
	jar.SetCookies(parse, c)
	client := &http.Client{
		Jar: jar,
	}
	return client, nil
}

func TestMiddlewareGzipHandleFunc(t *testing.T) {
	log, _ := logger.NewZapLogger(false)
	strg, _ := storage.NewStorage(log.Logger, "", "")
	hndlrs := handlers.NewHandlers(strg, log.Logger, "http://localhost:8080")

	adapter := HandlerFuncAdapter(hndlrs.PostRequestAPIShorten)
	//newHandler := MiddlewareGzipHandleFunc(nil)
	//handler := newHandler(adapter)

	server := httptest.NewServer(cookie.MiddlewareCheckCoockie(log.Logger,
		hndlrs)(MiddlewareGzipHandleFunc(log.Logger)(adapter)))
	defer server.Close()
	client, err := createHTTPAuthClient(log.Logger, server)
	assert.NoError(t, err)

	requestBody := `{"url": "http://test.com"}`

	t.Run("send_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		assert.NoError(t, err)
		err = zb.Close()
		assert.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, server.URL+"/api/shorten", buf)
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		request.Header.Set("Content-Encoding", "gzip")
		request.RequestURI = ""

		resp, err := client.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.NotNil(t, string(b))
	})

	t.Run("accept_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		request := httptest.NewRequest("POST", server.URL+"/api/shorten", buf)
		request.RequestURI = ""
		request.Header.Set("Accept-Encoding", "gzip")
		request.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := client.Do(request)
		assert.NoError(t, err)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		assert.NoError(t, err)

		b, err := io.ReadAll(zr)
		assert.NoError(t, err)
		assert.NotNil(t, string(b))
	})
}
