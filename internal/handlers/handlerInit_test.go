package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"practicumserver/internal/logger"
	"practicumserver/internal/storage"
	"strings"
	"testing"
)

// Функция тестирования Post запроса
func TestPostRequest(t *testing.T) {
	log, _ := logger.NewZapLogger(false)
	strg, _ := storage.NewStorage(log.Logger, "", "")
	hndlrs := NewHandlers(strg, log.Logger, "http://localhost:8080")

	type wantPost struct {
		statusCode  int
		contentType string
	}
	type request struct {
		body        io.Reader
		contentType string
	}
	tests := []struct {
		name    string
		want    wantPost
		url     string
		request request
	}{
		{
			name: "test one!",
			want: wantPost{
				statusCode:  201,
				contentType: "text/plain",
			},
			request: request{
				body:        strings.NewReader("http://hlijutdqqmefpt.net/zeosh/sthbp"),
				contentType: "text/plain; charset=utf-8",
			},
			url: "/",
		},
		{
			name: "test two!",
			want: wantPost{
				statusCode:  400,
				contentType: "",
			},
			request: request{
				body:        strings.NewReader("http://hlijutdqqmefpt.net/zeosh/sthbp"),
				contentType: "json",
			},
			url: "/",
		},
		{
			name: "test three!",
			want: wantPost{
				statusCode:  400,
				contentType: "",
			},
			request: request{
				body:        strings.NewReader(""),
				contentType: "text/plain; charset=utf-8",
			},
			url: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, tt.request.body)
			request.Header.Set("Content-Type", tt.request.contentType)
			w := httptest.NewRecorder()
			hndlrs.PostRequest(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			userResult, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)
			if res.StatusCode == http.StatusBadRequest {
				assert.Empty(t, userResult)
			} else if res.StatusCode == http.StatusCreated {
				assert.NotEmpty(t, userResult)
			}
		})
	}
}

// Функция тестирования Get запроса
func TestGetRequest(t *testing.T) {
	log, _ := logger.NewZapLogger(false)
	strg, _ := storage.NewStorage(log.Logger, "", "")
	hndlrs := NewHandlers(strg, log.Logger, "http://localhost:8080")

	type wantGet struct {
		statusCode int
		Location   string
	}
	tests := []struct {
		name   string
		want   wantGet
		adress string
	}{
		{
			name: "test one!",
			want: wantGet{
				statusCode: 307,
				Location:   "http://test",
			},
			adress: "/test",
		},
		{
			name: "test one!",
			want: wantGet{
				statusCode: 400,
				Location:   "",
			},
			adress: "/word",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.adress, nil)
			w := httptest.NewRecorder()
			hndlrs.GetRequest(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.Location, res.Header.Get("location"))
			err := res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestPostRequestApiShorten(t *testing.T) {
	log, _ := logger.NewZapLogger(false)
	strg, _ := storage.NewStorage(log.Logger, "", "")
	hndlrs := NewHandlers(strg, log.Logger, "http://localhost:8080")

	type wantPost struct {
		expectedCode int
		expectedBody string
		expectCt     string
	}
	type request struct {
		method      string
		body        string
		contentType string
		url         string
	}
	tests := []struct {
		name string
		request
		wantPost
	}{
		{
			name: "post request without body, result: status code 500",
			request: request{
				method:      http.MethodPost,
				body:        "",
				contentType: "application/json; charset=utf-8",
				url:         "http://localhost:8080/api/shorten",
			},
			wantPost: wantPost{
				expectedCode: http.StatusInternalServerError,
				expectedBody: "",
				expectCt:     "",
			},
		},
		{
			name: "method_post_with_wrong_content_type",
			request: request{
				method:      http.MethodPost,
				body:        "",
				contentType: "text/plain; charset=utf-8",
				url:         "http://localhost:8080/api/shorten",
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectedBody: "",
				expectCt:     "",
			},
		},
		{
			name: "method_post_with_any_url",
			request: request{
				method:      http.MethodPost,
				body:        "",
				contentType: "text/plain; charset=utf-8",
				url:         "http://localhost:8080/api/shorten/test",
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectedBody: "",
				expectCt:     "",
			},
		},
		{
			name: "method_post_unsupported_type_value",
			request: request{
				method:      http.MethodPost,
				body:        `{"url": 1}`,
				contentType: "application/json; charset=utf-8",
				url:         "http://localhost:8080/api/shorten",
			},
			wantPost: wantPost{
				expectedCode: http.StatusInternalServerError,
				expectedBody: "",
				expectCt:     "",
			},
		},
		{
			name: "method_post_success",
			request: request{
				method:      http.MethodPost,
				body:        `{"url":"http://test"}`,
				contentType: "application/json",
				url:         "http://localhost:8080/api/shorten",
			},
			wantPost: wantPost{
				expectedCode: http.StatusCreated,
				expectedBody: `{"result":"http://localhost:8080/test"}`,
				expectCt:     "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/shorten", strings.NewReader(tt.request.body))
			request.Header.Set("Content-Type", tt.request.contentType)
			w := httptest.NewRecorder()
			hndlrs.PostRequestAPIShorten(w, request)

			res := w.Result()

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.wantPost.expectedCode, res.StatusCode)
			assert.Equal(t, tt.wantPost.expectCt, res.Header.Get("Content-Type"))
			if tt.wantPost.expectedBody != "" {
				assert.JSONEq(t, tt.wantPost.expectedBody, string(resBody))
			}
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}