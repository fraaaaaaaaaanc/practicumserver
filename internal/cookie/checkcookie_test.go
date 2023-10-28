package cookie

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	handlers "practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
	"testing"
)

// createHTTPAuthClient creates an HTTP client with a JWT cookie for authentication in test func.
func createHTTPAuthClient(log *zap.Logger, srv *httptest.Server, userID string) *http.Client {
	// Create a cookie jar and JWT token for authentication.
	// Set up an HTTP client with the authenticated cookie.
	jar, _ := cookiejar.New(nil)
	token, err := BuildJWTString(userID)
	if err != nil {
		log.Error("Error test:", zap.Error(err))
		return nil
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
		log.Error("Error test:", zap.Error(err))
		return nil
	}
	jar.SetCookies(parse, c)
	client := &http.Client{
		Jar: jar,
	}
	return client
}

// TestMiddlewareCheckCookie is a test function for the MiddlewareCheckCookie middleware.
func TestMiddlewareCheckCookie(t *testing.T) {
	// Initialize logger, storage, and handlers for testing.
	log, _ := logger.NewZapLogger(false)
	strg, _ := storage.NewStorage(log.Logger, "", "")
	hndlrs := handlers.NewHandlers(strg, log.Logger, "http://localhost:8080")

	// Create an adapter for the handler function.
	adapter := models.HandlerFuncAdapter(hndlrs.GetRequestAPIUserUrls)

	// Create an HTTP server with middleware for GZIP compression.
	server := httptest.NewServer(MiddlewareCheckCookie(log.Logger,
		hndlrs)(adapter))
	defer server.Close()

	type req struct {
		client *http.Client
	}
	type resp struct {
		statusCode int
		respBody   string
	}

	tests := []struct {
		req
		resp
		userId string
		name   string
	}{
		{
			name: "test sending a request to the adress \"http://localhost:8080/api/user/urls\" " +
				"with userID = \"\", should return StatusCode: 401 and empty Response body",
			req: req{
				client: createHTTPAuthClient(log.Logger, server, ""),
			},
			resp: resp{
				statusCode: 401,
				respBody:   "",
			},
		},
		{
			name: "test sending a request to the adress \"http://localhost:8080/api/user/urls\" " +
				"with client: http.DefaultClient, should return StatusCode: 204 and empty Response body",
			req: req{
				client: http.DefaultClient,
			},
			resp: resp{
				statusCode: 204,
				respBody:   "",
			},
		},
		{
			name: "test sending a request to the adress \"http://localhost:8080/api/user/urls\" " +
				"with userID: test, should return StatusCode: 200 and is not empty Response body",
			req: req{
				client: createHTTPAuthClient(log.Logger, server, "test"),
			},
			resp: resp{
				statusCode: 200,
				respBody:   "[{\"short_url\":\"http://localhost:8080/test\",\"original_url\":\"http://test\"}]\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, server.URL+"/api/user/urls", nil)
			request.RequestURI = ""

			resp, err := tt.client.Do(request)
			assert.NoError(t, err)

			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, tt.statusCode, resp.StatusCode)
			assert.Equal(t, tt.respBody, string(respBody))
		})
	}
}
