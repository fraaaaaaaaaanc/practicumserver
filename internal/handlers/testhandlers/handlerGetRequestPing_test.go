package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	handlers "practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/storage/pg"
	"testing"
)

func TestGetRequestPing(t *testing.T) {
	log, _ := logger.NewZapLogger(false)

	newCookie := &http.Cookie{
		Name: "Authorization",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTgxMTI0ODksIlVzZXJJRCI6IktacjlENG01bkJuY05uMUNQ" +
			"M08xbHc9PSJ9.g0vISaj4K1rP4V83AOD8Q4y4_0gsZ6Dwci1eZ72jM54",
		Path:     "/",
		MaxAge:   7200,
		HttpOnly: true,
	}

	type requestData struct {
		dbAdress string
		adress   string
	}
	tests := []struct {
		name       string
		statusCode int
		cookie     *http.Cookie
		requestData
	}{
		//{
		//	name: "test sending a request to the address \"http://localhost:8080/ping\", while initializing " +
		//		"NewHandlers with the correct dbadress parameter, should return Status Code 200",
		//	statusCode: http.StatusOK,
		//	requestData: requestData{
		//		dbAdress: "host=localhost user=postgres password=1234 dbname=linksShorten sslmode=disable",
		//		adress:   "/ping",
		//	},
		//},
		//{
		//	name: "test sending a request to the address \"http://localhost:8080/ping\", while initializing " +
		//		"NewHandlers with an incorrect dbadress parameter, should return Status Code 200",
		//	statusCode: http.StatusBadRequest,
		//	requestData: requestData{
		//		dbAdress: "host=localhost user=postgres dbname=linksShorten sslmode=disable",
		//		adress:   "/ping",
		//	},
		//},
		{
			name: "test sending a request to the address \"http://localhost:8080/ping\", while initializing " +
				"NewHandlers with an empty dbadress parameter, should return Status Code 200",
			statusCode: http.StatusBadRequest,
			cookie:     newCookie,
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
			err := res.Body.Close()
			require.NoError(t, err)
		})
	}
}
