package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"practicumserver/internal/models"
)

// PostRequestAPIShortenBatch is an HTTP handler method that processes POST requests to the "/api/shorten/batch" endpoint.
// It handles batch requests with JSON data in the request body, where each item in the JSON array represents a URL to be shortened.
// This handler performs the following steps:
// 1. Checks the request URL to ensure it matches the expected endpoint ("/api/shorten/batch").
// 2. Parses and decodes JSON data from the request body into a slice of models.RequestAPIBatch structures.
// 3. Validates each original URL provided in the request data.
// 4. Calls the Storage.SetListData method to store the original URLs and generate short links for each.
// 5. Responds with the appropriate HTTP status code and a JSON response containing the short links corresponding to the original URLs.
//
// This handler is designed to process multiple URLs in a single batch request efficiently.
func (h *Handlers) PostRequestAPIShortenBatch(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/api/shorten/batch" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL"))
		return
	}

	var req []models.RequestAPIBatch
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	for _, originalURLStruct := range req {
		if _, err := url.ParseRequestURI(originalURLStruct.OriginalURL); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.Log.Error("Error:",
				zap.String("reason", "The request body isn't a url"))
			return
		}
	}

	resp, err := h.Storage.SetListData(r.Context(), req, h.prefix)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err = enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
}
