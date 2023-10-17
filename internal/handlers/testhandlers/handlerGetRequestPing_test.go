package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	handlers "practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/storage"
	"testing"
)

func TestGetRequestPing(t *testing.T) {
	log, _ := logger.NewZapLogger(false)
	type requestData struct {
		dbAdress string
		adress   string
	}
	tests := []struct {
		name       string
		statusCode int
		requestData
	}{
		{
			name: "test sending a request to the address \"http://localhost:8080/ping\", while initializing " +
				"NewHandlers with the correct dbadress parameter, should return Status Code 200",
			statusCode: http.StatusOK,
			requestData: requestData{
				dbAdress: "host=localhost user=postgres password=1234 dbname=linksShorten sslmode=disable",
				adress:   "/ping",
			},
		},
		{
			name: "test sending a request to the address \"http://localhost:8080/ping\", while initializing " +
				"NewHandlers with an incorrect dbadress parameter, should return Status Code 200",
			statusCode: http.StatusBadRequest,
			requestData: requestData{
				dbAdress: "host=localhost user=postgres dbname=linksShorten sslmode=disable",
				adress:   "/ping",
			},
		},
		{
			name: "test sending a request to the address \"http://localhost:8080/ping\", while initializing " +
				"NewHandlers with an empty dbadress parameter, should return Status Code 200",
			statusCode: http.StatusBadRequest,
			requestData: requestData{
				dbAdress: "",
				adress:   "/ping",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strg, _ := storage.NewStorage(log.Logger, tt.dbAdress, "")
			hndlrs := handlers.NewHandlers(strg, log.Logger, "http://localhost:8080")

			request := httptest.NewRequest(http.MethodGet, tt.adress, nil)
			w := httptest.NewRecorder()
			hndlrs.GetRequestPing(w, request)

			res := w.Result()
			assert.Equal(t, tt.statusCode, res.StatusCode)
		})
	}
}
