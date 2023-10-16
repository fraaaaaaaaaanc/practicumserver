package handlers

import (
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"practicumserver/internal/utils"
)

func (h *Handlers) PostRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if (!utils.ValidContentType(contentType, "text/plain") && contentType != "application/x-gzip") ||
		r.URL.String() != "/" {
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(h.prefix + "/" + newShortLink))
}
