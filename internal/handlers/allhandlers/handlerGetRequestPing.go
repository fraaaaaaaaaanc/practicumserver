package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/storage"
)

func (h *Handlers) GetRequestPing(w http.ResponseWriter, r *http.Request) {
	if ds, ok := h.Storage.(*storage.DBStorage); ok {
		if err := ds.PingDB(r.Context()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.Log.Error("Error:", zap.Error(err))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusBadRequest)
}
