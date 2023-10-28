package handlers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
	"strings"
	"testing"
)

// TestPostRequestApiShorten tests the function PostRequestApiShorten
func TestPostRequestApiShorten(t *testing.T) {
	log, _ := logger.NewZapLogger(false)
	strg, _ := storage.NewStorage(log.Logger, "", "")
	hndlrs := handlers.NewHandlers(strg, log.Logger, "http://localhost:8080")

	newCookie := &http.Cookie{
		Name: "Authorization",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTgxMTI0ODksIlVzZXJJRCI6IktacjlENG01bkJuY05uMUNQ" +
			"M08xbHc9PSJ9.g0vISaj4K1rP4V83AOD8Q4y4_0gsZ6Dwci1eZ72jM54",
		Path:     "/",
		MaxAge:   7200,
		HttpOnly: true,
	}

	type wantPost struct {
		expectedCode int
		expectCt     string
	}
	type request struct {
		method      string
		body        string
		contentType string
		url         string
		cookie      *http.Cookie
	}
	tests := []struct {
		name string
		request
		wantPost
	}{
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten\" with an empty request body, it " +
				"should return the Status Code 400",
			request: request{
				method:      http.MethodPost,
				body:        "",
				contentType: "application/json; charset=utf-8",
				url:         "http://localhost:8080/api/shorten",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectCt:     "",
			},
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten\" with the request body in the  " +
				"appropriate json format, should return the Status Code 400",
			request: request{
				method:      http.MethodPost,
				body:        `http:localhost"`,
				contentType: "application/json; charset=utf-8",
				url:         "http://localhost:8080/api/shorten",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectCt:     "",
			},
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten\" with a request body in which " +
				"the parameter does not match the models structure type models.RequestAPIShorten, should return " +
				"the Status Code 400",
			request: request{
				method:      http.MethodPost,
				body:        `{"url": 1}`,
				contentType: "application/json; charset=utf-8",
				url:         "http://localhost:8080/api/shorten",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectCt:     "",
			},
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten\" with a request body in which the " +
				"parameter is not a URL, the Status Code 400",
			request: request{
				method:      http.MethodPost,
				body:        `{"url": Hello}`,
				contentType: "application/json; charset=utf-8",
				url:         "http://localhost:8080/api/shorten",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectCt:     "",
			},
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten\" with the correct request body," +
				"return the Status Code 201",
			request: request{
				method:      http.MethodPost,
				body:        `{"url":"http://newTest"}`,
				contentType: "application/json",
				url:         "http://localhost:8080/api/shorten",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusCreated,
				expectCt:     "application/json",
			},
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten\" with the request body as in the previous" +
				"request, return the Status Code 409",
			request: request{
				method:      http.MethodPost,
				body:        `{"url":"http://newTest"}`,
				contentType: "application/json",
				url:         "http://localhost:8080/api/shorten",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusConflict,
				expectCt:     "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/shorten", strings.NewReader(tt.request.body))
			request.Header.Set("Content-Type", tt.request.contentType)
			ctx := context.WithValue(request.Context(), models.UserIDKey, "1234")
			request = request.WithContext(ctx)
			w := httptest.NewRecorder()
			hndlrs.PostRequestAPIShorten(w, request)

			res := w.Result()

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.wantPost.expectedCode, res.StatusCode)
			assert.Equal(t, tt.wantPost.expectCt, res.Header.Get("Content-Type"))
			if tt.request.body != "" {
				assert.NotNil(t, resBody)
			}
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}