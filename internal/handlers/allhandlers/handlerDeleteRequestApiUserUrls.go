package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/models"
)

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
