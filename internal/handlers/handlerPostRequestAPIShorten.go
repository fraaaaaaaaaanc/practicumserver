package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"practicumserver/internal/models"
	"practicumserver/internal/utils"
)

func (h *Handlers) PostRequestAPIShorten(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if (!utils.ValidContentType(contentType, "application/json") && contentType != "application/x-gzip") ||
		r.URL.String() != "/api/shorten" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL or Content-Type"))
		return
	}

	var req models.RequestApiShorten
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	if _, err := url.ParseRequestURI(req.LongURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "The request body isn't a url"))
		return
	}

	newShortLink, err := h.Storage.SetData(r.Context(), req.LongURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	resp := models.ResponseApiShorten{
		ShortURL: h.prefix + "/" + newShortLink,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
}
