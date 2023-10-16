package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

func (h *Handlers) GetRequest(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.String()[1:]
	baseLink, err := h.Storage.GetData(r.Context(), shortLink)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
	}
	if baseLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error: This abbreviated link wasn't found",
			zap.String("shortLink", shortLink))
		return
	}
	w.Header().Set("Location", baseLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
