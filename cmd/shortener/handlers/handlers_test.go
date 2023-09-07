package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type wantPost struct {
	statusCode  int
	contentType string
}
type request struct {
	body        io.Reader
	contentType string
}

// Функция тестирования Post запроса
func TestPostRequest(t *testing.T) {
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
			PostRequest(w, request)

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
			adress: "//test?id=test",
		},
		{
			name: "test two!",
			want: wantGet{
				statusCode: 400,
				Location:   "",
			},
			adress: "/word",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(Router())
			defer ts.Close()
			req, err := http.NewRequest(http.MethodGet, ts.URL+tt.adress, nil)
			res, err := ts.Client().Do(req)
			require.NoError(t, err)
			//w := httptest.NewRecorder()
			//GetRequest(w, request)

			//res := w.Result()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.Location, res.Header.Get("location"))
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}
