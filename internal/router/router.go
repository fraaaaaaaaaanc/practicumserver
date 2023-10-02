package router

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/compress"
	"practicumserver/internal/config"
	"practicumserver/internal/handlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/storage"
)

func Router(flags *config.Flags, storage *storage.Storage, log *zap.Logger) chi.Router {
	var handlers handlers.Handlers

	r := chi.NewRouter()

	r.Use(logger.MiddlewareLogHandleFunc(log), compress.MiddlewareGzipHandleFunc)
	r.Get("/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetRequest(w, r, storage)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostRequest(w, r, storage, flags.ShortLink)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostRequestAPIShorten(w, r, storage, flags.ShortLink)
	})

	return r
}
