package router

import (
	"github.com/go-chi/chi"
	"net/http"
	"practicumserver/cmd/shortener/config"
	"practicumserver/cmd/shortener/handlers"
)

func Router(flags *config.Flags) chi.Router {
	r := chi.NewRouter()
	r.Get("/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetRequest(w, r)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostRequest(w, r, flags)
	})

	return r
}
