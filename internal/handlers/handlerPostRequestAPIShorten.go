package handlers

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
)

func (h *Handlers) PostRequestAPIShorten(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/api/shorten" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL or Content-Type"))
		return
	}

	var req models.RequestAPIShorten
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
	if err != nil && !errors.Is(err, storage.ErrConflictData) {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	resp := models.ResponseAPIShorten{
		ShortURL: h.prefix + "/" + newShortLink,
	}
	w.Header().Set("Content-Type", "application/json")
	httpStatus := http.StatusCreated
	if errors.Is(err, storage.ErrConflictData) {
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
