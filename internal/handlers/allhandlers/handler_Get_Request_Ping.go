package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/storage/pg"
)

// PingRequest is an HTTP handler method that processes GET requests to the "/ping" endpoint.
// It is a simple health check endpoint to verify that the server is up and running.
// This handler responds with an HTTP 200 OK status code and the message "pong" as a simple acknowledgment.
func (h *Handlers) GetRequestPing(w http.ResponseWriter, r *http.Request) {
	if ds, ok := h.Storage.(*pgstorage.DBStorage); ok {
		if err := ds.PingDB(r.Context()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.Log.Error("Error:", zap.Error(err))
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
