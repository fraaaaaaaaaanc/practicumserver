package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handlers) GetRequestAPIUserUrls(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("userID") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		h.Log.Info("Info: Cookie do not have userID")
		return
	}

	resp, err := h.Storage.GetListData(r.Context(), h.prefix)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	if len(resp) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err = enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
}
