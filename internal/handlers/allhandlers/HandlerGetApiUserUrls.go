package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handlers) GetRequestApiUserUrls(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Storage.GetListData(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
	if len(resp) == 0 {
		w.WriteHeader(http.StatusNoContent)
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
