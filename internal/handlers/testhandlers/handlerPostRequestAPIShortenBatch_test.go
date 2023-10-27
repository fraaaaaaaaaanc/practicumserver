package handlers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	handlers "practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
	"strings"
	"testing"
)

func TestPostRequestAPIShortenBatch(t *testing.T) {
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
		expectedBody string
		expectCt     string
	}
	type request struct {
		method      string
		body        string
		contentType string
		cookie      *http.Cookie
	}
	tests := []struct {
		url  string
		name string
		request
		wantPost
	}{
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten/batch\" with an empty request body, " +
				"it should return the Status Code 400",
			request: request{
				method:      http.MethodPost,
				body:        "",
				contentType: "application/json",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectCt:     "",
			},
			url: "/api/shorten/batch",
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten/batch\" with the request body in the  " +
				"appropriate json format, should return the Status Code 400",
			request: request{
				method:      http.MethodPost,
				body:        "http://localhost",
				contentType: "application/json",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectCt:     "",
			},
			url: "/api/shorten/batch",
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten/batch\" with a request body in which " +
				"the first parameter does not match the models structure type models.RequestAPIShorten, should return " +
				"the Status Code 400",
			request: request{
				method:      http.MethodPost,
				body:        `{"correlation_id": 1,"original_url": "http://ya.com"}`,
				contentType: "application/json",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectCt:     "",
			},
			url: "/api/shorten/batch",
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten/batch\" with a request body in which " +
				"the two parameter does not match the models structure type models.RequestAPIShorten, should return " +
				"the Status Code 400",
			request: request{
				method:      http.MethodPost,
				body:        `{"correlation_id": "1","original_url": 2}`,
				contentType: "application/json",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectCt:     "",
			},
			url: "/api/shorten/batch",
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten/batch\" with a request body in " +
				"which the parameter is not a URL, the Status Code 400",
			request: request{
				method:      http.MethodPost,
				body:        `{"correlation_id": "1","original_url": Hello}`,
				contentType: "application/json",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusBadRequest,
				expectCt:     "",
			},
			url: "/api/shorten/batch",
		},
		{
			name: "POST request was sent to \"http://localhost:8080/api/shorten/batch\" with the correct request" +
				"body, return the Status Code 201",
			request: request{
				method:      http.MethodPost,
				body:        `[{"correlation_id": "1","original_url": "http://newTest"},{"correlation_id": "2","original_url": "http://testNew"}]`,
				contentType: "application/json",
				cookie:      newCookie,
			},
			wantPost: wantPost{
				expectedCode: http.StatusCreated,
				expectCt:     "application/json",
			},
			url: "/api/shorten/batch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.url, strings.NewReader(tt.request.body))
			request.Header.Set("Content-Type", tt.request.contentType)
			ctx := context.WithValue(request.Context(), models.UserIDKey, "1234")
			request = request.WithContext(ctx)
			w := httptest.NewRecorder()
			hndlrs.PostRequestAPIShortenBatch(w, request)

			resp := w.Result()
			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantPost.expectedCode, resp.StatusCode)
			assert.Equal(t, tt.wantPost.expectCt, resp.Header.Get("Content-Type"))
			if tt.request.body != "" {
				assert.NotNil(t, respBody)
			}
			err = resp.Body.Close()
			assert.NoError(t, err)
		})
	}
}
