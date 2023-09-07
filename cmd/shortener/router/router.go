package router

import (
	"github.com/go-chi/chi"
	"practicumserver/cmd/shortener/handlers"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id:[a-zA-Z0-9]+}", handlers.GetRequest)
	r.Post("/", handlers.PostRequest)

	return r
}
