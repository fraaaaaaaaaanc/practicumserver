package handlers

import (
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"practicumserver/internal/models"
)

// PostRequest is an HTTP handler method that processes POST requests to the root ("/") endpoint.
// It handles requests with the Content-Type set to text/plain and expects the request body to contain a single URL.
// This handler performs the following steps:
// 1. Checks the request URL to ensure it matches the root endpoint ("/").
// 2. Reads the request body to obtain the original URL.
// 3. Validates the provided URL.
// 4. Calls the Storage.SetData method to store the original URL and generate a short link.
// 5. Responds with the appropriate HTTP status code and the short link if the URL was successfully shortened.
//
// This handler is designed to process individual URL shortening requests with text/plain content.
func (h *Handlers) PostRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL or Content-Type"))
		return
	}

	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	if _, err := url.ParseRequestURI(string(originalURL)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "The request body isn't a url"))
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			//w.WriteHeader(http.StatusInternalServerError)
			//h.Log.Error("Error:", zap.Error(err))
			return
		}
	}()

	shortLink, err := h.Storage.SetData(r.Context(), string(originalURL))
	if err != nil && !errors.Is(err, models.ErrConflictData) {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	httpStatus := http.StatusCreated
	if errors.Is(err, models.ErrConflictData) {
		httpStatus = http.StatusConflict
	}
	w.WriteHeader(httpStatus)
	_, _ = w.Write([]byte(h.prefix + "/" + shortLink))
}
