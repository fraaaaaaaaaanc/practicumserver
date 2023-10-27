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

// Функция тестирования Post запроса
func TestPostRequest(t *testing.T) {
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
		statusCode  int
		contentType string
	}
	type request struct {
		body        string
		contentType string
		cookie      *http.Cookie
	}
	tests := []struct {
		name    string
		want    wantPost
		url     string
		request request
	}{
		{
			name: "POST request was sent to \"http://localhost:8080/\" with an empty request body, it should " +
				"return the Status Code 400",
			want: wantPost{
				statusCode:  http.StatusBadRequest,
				contentType: "",
			},
			request: request{
				body:        "",
				contentType: "text/plain; charset=utf-8",
				cookie:      newCookie,
			},
			url: "/",
		},
		{
			name: "a POST request was sent to \"http://localhost:8080/\" with a request body that is not a url," +
				"return the Status Code 400",
			want: wantPost{
				statusCode:  http.StatusBadRequest,
				contentType: "",
			},
			request: request{
				body:        "notLink",
				contentType: "text/plain; charset=utf-8",
				cookie:      newCookie,
			},
			url: "/",
		},
		{
			name: "a POST request was sent to \"http://localhost:8080/\" with the correct request body," +
				"return the Status Code 400",
			want: wantPost{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
			},
			request: request{
				body:        "http://newTest",
				contentType: "text/plain; charset=utf-8",
				cookie:      newCookie,
			},
			url: "/",
		},
		{
			name: "a POST request was sent to \"http://localhost:8080/\" with the same request body as in the last " +
				"test, return the Status Code 409",
			want: wantPost{
				statusCode:  http.StatusConflict,
				contentType: "text/plain",
			},
			request: request{
				body:        "http://newTest",
				contentType: "text/plain; charset=utf-8",
				cookie:      newCookie,
			},
			url: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, strings.NewReader(tt.request.body))
			request.Header.Set("Content-Type", tt.request.contentType)
			ctx := context.WithValue(request.Context(), models.UserIDKey, "1234")
			request = request.WithContext(ctx)
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
