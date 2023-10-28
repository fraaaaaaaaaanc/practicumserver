package handlers

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"practicumserver/internal/models"
)

// PostRequestAPIShorten is an HTTP handler method that processes POST requests to the "/api/shorten" endpoint.
// It handles requests with JSON data in the request body. This handler performs the following steps:
// 1. Checks the request URL to ensure it matches the expected endpoint ("/api/shorten").
// 2. Parses and decodes JSON data from the request body into a models.RequestAPIShorten structure.
// 3. Validates the original URL provided in the request body.
// 4. Calls the Storage.SetData method to store the original URL and generate a new short link.
// 5. Responds with the appropriate HTTP status code and a JSON response containing the short link.
//
// If the request contains data that is already stored in the storage (models.ErrConflictData), it returns a
// HTTP status 409 (Conflict), indicating that the resource already exists.
func (h *Handlers) PostRequestAPIShorten(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/api/shorten" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL"))
		return
	}

	var req models.RequestAPIShorten
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	if _, err := url.ParseRequestURI(req.OriginalURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "The request body isn't a url"))
		return
	}

	newShortLink, err := h.Storage.SetData(r.Context(), req.OriginalURL)
	if err != nil && !errors.Is(err, models.ErrConflictData) {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	resp := models.ResponseAPIShorten{
		ShortLink: h.prefix + "/" + newShortLink,
	}
	w.Header().Set("Content-Type", "application/json")
	httpStatus := http.StatusCreated
	if errors.Is(err, models.ErrConflictData) {
		httpStatus = http.StatusConflict
	}
	w.WriteHeader(httpStatus)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
}
