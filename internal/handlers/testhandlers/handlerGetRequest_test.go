package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/storage"
	"testing"
)

// Функция тестирования Get запроса
func TestGetRequest(t *testing.T) {
	log, _ := logger.NewZapLogger(false)
	strg, _ := storage.NewStorage(log.Logger, "", "")
	hndlrs := handlers.NewHandlers(strg, log.Logger, "http://localhost:8080")

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
			name: "test sending a request to the address \"http://localhost:8080 /\" , should return Status Code 400",
			want: wantGet{
				statusCode: http.StatusBadRequest,
				Location:   "",
			},
			adress: "/",
		},
		{
			name: "test sending a request to the address \"http://localhost:8080/word \" despite the fact that the " +
				"abbreviated word link is not written to the repository, it should return the Status Code 400",
			want: wantGet{
				statusCode: http.StatusBadRequest,
				Location:   "",
			},
			adress: "/word",
		},
		{
			name: "test sending a request to the address \"http://localhost:8080/test \" while the abbreviated test " +
				"link is written to the repository, it should return Status Code 307 and location \"http://test \"",
			want: wantGet{
				statusCode: http.StatusTemporaryRedirect,
				Location:   "http://test",
			},
			adress: "/test",
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
