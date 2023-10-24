package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

// Хендлер для обратотки GET запросов по адресу "/shortLink"
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
	//Проверяем полученный оригинальый URL на пустоту
	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error: This abbreviated shortLink wasn't found",
			zap.String("shortLink", shortLink))
		return
	}
	//Устанавливаем Location и статус код StatusTemporaryRedirect
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
