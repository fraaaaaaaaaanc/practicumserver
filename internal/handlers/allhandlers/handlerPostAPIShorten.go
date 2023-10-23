package handlers

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
)

// Хендер принимающий POST запрос по адрессу "/api/shorten", в теле которого могут лежать данные в формате JSON
func (h *Handlers) PostRequestAPIShorten(w http.ResponseWriter, r *http.Request) {
	//Проверка адресса
	if r.URL.String() != "/api/shorten" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL"))
		return
	}

	//Перенос данных из JSON формата в структуру req
	var req models.RequestAPIShorten
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	//Проверка полученных данных на соотвествие URL
	if _, err := url.ParseRequestURI(req.OriginalURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "The request body isn't a url"))
		return
	}

	//Проверка полученных данных на соотвествие URL
	newShortLink, err := h.Storage.SetData(r.Context(), h.prefix, req.OriginalURL)
	if err != nil && !errors.Is(err, storage.ErrConflictData) {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	//Формирование ответа
	resp := models.ResponseAPIShorten{
		ShortLink: h.prefix + "/" + newShortLink,
	}
	w.Header().Set("Content-Type", "application/json")
	//Метод SetData может вернуть ошибку типа ErrConflictData, это означает что в запросе были
	//полученны данные которые уже записаны в хранилище, поэтому в таком случае выставдляется статус 409
	httpStatus := http.StatusCreated
	if errors.Is(err, storage.ErrConflictData) {
		httpStatus = http.StatusConflict
	}
	w.WriteHeader(httpStatus)

	//Перенос данных из структуры в формат JSON
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
}
