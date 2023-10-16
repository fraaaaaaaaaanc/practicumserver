package router

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/compress"
	"practicumserver/internal/handlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/storage"
)

func Router(log *zap.Logger, prefix string,
	dbStorageAdress, fileStoragePath string) (chi.Router, error) {
	strg, err := storage.NewStorage(log, dbStorageAdress, fileStoragePath)
	if err != nil {
		return nil, err
	}
	hndlrs := handlers.NewHandlers(strg, log, prefix)

	r := chi.NewRouter()

	r.Use(logger.MiddlewareLogHandleFunc(log), compress.MiddlewareGzipHandleFunc(log))
	r.Get("/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.GetRequest(w, r)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.GetRequestPing(w, r)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.PostRequest(w, r)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.PostRequestAPIShorten(w, r)
	})
	r.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.PostRequestAPIShortenBatch(w, r)
	})

	return r, nil
}
