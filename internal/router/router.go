// Package router provides functionality for setting up and configuring the application's routing.
// The Router function initializes and configures routes using the go-chi/chi router package. It sets up
// various endpoints for handling HTTP requests, along with the corresponding handlers defined in the
// handlers package.
package router

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/compress"
	coockie "practicumserver/internal/cookie"
	"practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
)

// Router creates and configures the application's router.
// It sets up various endpoints for handling HTTP requests, including:
// - Endpoint to retrieve the original URL from a shortened URL.
// - Endpoint for checking the connectivity to the database.
// - Endpoint to retrieve user URLs.
// - Endpoint to create a shortened URL.
// - Endpoint to create a shortened URL with JSON input.
// - Endpoint to create shortened URLs in batches.
// - Endpoint to delete user URLs.
// The provided handlers and logger are used to process requests and log relevant information.
func Router(hndlrs *handlers.Handlers, log *zap.Logger) (chi.Router, error) {
	r := chi.NewRouter()

	r.Use(coockie.MiddlewareCheckCookie(log, hndlrs), logger.MiddlewareLogHandleFunc(log),
		compress.MiddlewareGzipHandleFunc(log))
	r.Get("/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.GetRequest(w, r)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.GetRequestPing(w, r)
	})
	r.Get("/api/user/urls", hndlrs.GetRequestAPIUserUrls)
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.PostRequest(w, r)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.PostRequestAPIShorten(w, r)
	})
	r.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		hndlrs.PostRequestAPIShortenBatch(w, r)
	})
	r.Delete("/api/user/urls", hndlrs.DeleteRequestAPIUserUrls)

	return r, nil
}
