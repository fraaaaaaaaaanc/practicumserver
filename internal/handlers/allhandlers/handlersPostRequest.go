package handlers

import (
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"practicumserver/internal/storage"
)

func (h *Handlers) PostRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL or Content-Type"))
		return
	}
	link, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	if _, err := url.ParseRequestURI(string(link)); err != nil {
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

	newShortLink, err := h.Storage.SetData(r.Context(), string(link))
	if err != nil && !errors.Is(err, storage.ErrConflictData) {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	httpStatus := http.StatusCreated
	if errors.Is(err, storage.ErrConflictData) {
		httpStatus = http.StatusConflict
	}
	w.WriteHeader(httpStatus)
	_, _ = w.Write([]byte(h.prefix + "/" + newShortLink))
}
