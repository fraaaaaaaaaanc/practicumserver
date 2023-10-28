package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/models"
)

// DeleteRequestAPIUserUrls is an HTTP handler method that processes DELETE requests from path "/api/user/urls",
// to remove user-specific short URLs.
// It expects a JSON array of short URLs in the request body and triggers the deletion process for each URL.
// The user's authentication and authorization information is obtained from the request's context.
// If the request body is not a valid JSON array or an error occurs during processing, it returns a HTTP 400 Bad Request response.
// Upon successful processing, it responds with an HTTP 202 Accepted status code.
func (h *Handlers) DeleteRequestAPIUserUrls(w http.ResponseWriter, r *http.Request) {
	var URLList models.DeleteURLList
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&URLList); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	for _, shortLink := range URLList {
		deleteURL := &models.DeleteURL{
			UserID:   r.Context().Value(models.UserIDKey).(string),
			ShortURL: shortLink,
		}
		h.DelCn <- deleteURL
	}
	w.WriteHeader(http.StatusAccepted)
}
