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
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
	"testing"
)

// createHTTPAuthClient creates an HTTP client with a JWT cookie for authentication in test func.
func createHTTPAuthClient(log *zap.Logger, srv *httptest.Server) (*http.Client, error) {
	// Create a cookie jar and JWT token for authentication.
	// Set up an HTTP client with the authenticated cookie.
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

// TestMiddlewareGzipHandleFunc is a test function for the MiddlewareGzipHandleFunc middleware.
func TestMiddlewareGzipHandleFunc(t *testing.T) {
	// Initialize logger, storage, and handlers for testing.
	log, _ := logger.NewZapLogger(false)
	strg, _ := storage.NewStorage(log.Logger, "", "")
	hndlrs := handlers.NewHandlers(strg, log.Logger, "http://localhost:8080")

	// Create an adapter for the handler function.
	adapter := models.HandlerFuncAdapter(hndlrs.PostRequestAPIShorten)

	// Create an HTTP server with middleware for GZIP compression.
	server := httptest.NewServer(cookie.MiddlewareCheckCookie(log.Logger,
		hndlrs)(MiddlewareGzipHandleFunc(log.Logger)(adapter)))
	defer server.Close()
	client, err := createHTTPAuthClient(log.Logger, server)
	assert.NoError(t, err)

	requestBody := `{"url": "http://test.com"}`

	t.Run("send_gzip", func(t *testing.T) {
		// Test case for sending GZIP-encoded data.
		// Compress the request body and make an HTTP POST request with GZIP-encoded data.
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
		// Test case for accepting GZIP-encoded data.
		// Send an HTTP POST request with an "Accept-Encoding: gzip" header.
		buf := bytes.NewBufferString(requestBody)
		request := httptest.NewRequest(http.MethodPost, server.URL+"/api/shorten", buf)
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
