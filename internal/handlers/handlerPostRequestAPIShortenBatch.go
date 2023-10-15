package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/models"
)

func (h *Handlers) PostRequestAPIShortenBatch(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/api/shorten/batch" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL or Content-Type"))
		return
	}

	var req []models.RequestApiBatch
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	resp, err := h.Storage.SetListData(r.Context(), req)
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
