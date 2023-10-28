package handlers

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/models"
)

// GetRequest is an HTTP handler method that processes GET requests from path "/{shortLink}",
// to redirect users to the original URL corresponding to a short link.
// It extracts the short link from the request URL and checks if it's empty. If empty, it responds with an
// HTTP 400 Bad Request status code.
// It calls the GetData method to retrieve the original URL. If the original URL is not found, it returns an
// HTTP 400 Bad Request status code.
// If the short link points to a deleted URL, it responds with an HTTP 410 Gone status code.
// If the original URL is found, it sets the "Location" header to the original URL and responds with an
// HTTP 307 Temporary Redirect status code.
func (h *Handlers) GetRequest(w http.ResponseWriter, r *http.Request) {
	//Получем сокращенную ссылку из адреса и проверяем ее на пустоту
	shortLink := r.URL.String()[1:]
	if shortLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error: An empty link was sent to the get request")
		return
	}
	//Отправляем shortLink в метод GetData для поиска оригинального URL,
	//если оригинкальны URL не найден вощзвращается StatusBadRequest
	originalURL, err := h.Storage.GetData(r.Context(), shortLink)
	if err != nil && !errors.Is(err, models.ErrDeletedData) {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	if errors.Is(err, models.ErrDeletedData) {
		w.WriteHeader(http.StatusGone)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	//Устанавливаем Location и статус код StatusTemporaryRedirect
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
