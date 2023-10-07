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

func Router(log *zap.Logger, prefix, fileStorage, dbAdress string) chi.Router {
	hndlrs := handlers.NewHandlers(prefix, dbAdress, fileStorage)
	storage.NewRead(fileStorage, hndlrs.Storage)

	r := chi.NewRouter()

	r.Use(logger.MiddlewareLogHandleFunc(log), compress.MiddlewareGzipHandleFunc)
	r.Get("/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.GetRequest(w, r)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.GerRequestPing(w, r)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.PostRequest(w, r)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.PostRequestAPIShorten(w, r)
	})

	return r
}
