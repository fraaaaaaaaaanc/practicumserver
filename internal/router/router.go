package router

import (
	"github.com/go-chi/chi"
	"net/http"
	"practicumserver/internal/config"
	handlers2 "practicumserver/internal/handlers"
	"practicumserver/internal/logger"
	storage2 "practicumserver/internal/storage"
)

func Router(flags *config.Flags, storage *storage2.Storage) chi.Router {
	var handlers handlers2.Handlers

	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Get("/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetRequest(w, r, storage)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostRequest(w, r, storage, flags)
	})

	return r
}
