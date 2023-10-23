package handlers

import (
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"practicumserver/internal/storage"
)

// Хендер принимающий POST запрос по адрессу "/"
func (h *Handlers) PostRequest(w http.ResponseWriter, r *http.Request) {
	//Проверка адресса
	if r.URL.String() != "/" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL or Content-Type"))
		return
	}
	//Считывание тела запроса
	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	//Проверка полученных данных на соотвествие URL
	if _, err := url.ParseRequestURI(string(originalURL)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "The request body isn't a url"))
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			//w.WriteHeader(http.StatusInternalServerError)
			//h.Log.Error("Error:", zap.Error(err))
			return
		}
	}()

	//Проверка полученных данных на соотвествие URL
	shortLink, err := h.Storage.SetData(r.Context(), h.prefix, string(originalURL))
	if err != nil && !errors.Is(err, storage.ErrConflictData) {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	//Формирование ответа
	w.Header().Set("Content-Type", "text/plain")
	//Метод SetData может вернуть ошибку типа ErrConflictData, это означает что в запросе были
	//полученны данные которые уже записаны в хранилище, поэтому в таком случае выставдляется статус 409
	httpStatus := http.StatusCreated
	if errors.Is(err, storage.ErrConflictData) {
		httpStatus = http.StatusConflict
	}
	w.WriteHeader(httpStatus)
	_, _ = w.Write([]byte(shortLink))
}
