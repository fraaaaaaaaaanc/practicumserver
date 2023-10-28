package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/models"
)

// GetRequestAPIUserUrls is an HTTP handler method that processes GET requests from path "/api/user/urls",
// to retrieve a list of user-specific short URLs.
// It checks the user's authentication and authorization using the "userID" value from the request context.
// If the user is not authenticated, it responds with an HTTP 401 Unauthorized status code.
// It fetches the list of short URLs from the storage component and responds with them in JSON format.
// If the retrieval process encounters an error, it returns an HTTP 400 Bad Request status code.
// If no short URLs are found for the user, it responds with an HTTP 204 No Content status code.
func (h *Handlers) GetRequestAPIUserUrls(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(models.UserIDKey).(string) == "" {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err = enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
}
