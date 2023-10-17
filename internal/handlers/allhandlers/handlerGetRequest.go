package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

func (h *Handlers) GetRequest(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.String()[1:]
	if shortLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error: An empty link was sent to the get request")
		return
	}
	baseLink, err := h.Storage.GetData(r.Context(), shortLink)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
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
