package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"practicumserver/internal/models"
)

// Хендер принимающий POST запрос по адрессу "/api/shorten/Batch", в теле которого может лежать пачка данных в формате JSON
func (h *Handlers) PostRequestAPIShortenBatch(w http.ResponseWriter, r *http.Request) {
	//Проверка адресса
	if r.URL.String() != "/api/shorten/batch" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL"))
		return
	}

	//Перенос данных из JSON формата в cлайс структур req
	var req []models.RequestAPIBatch
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	//Проверка полученных данных на соотвествие URL
	for _, originalURLStruct := range req {
		if _, err := url.ParseRequestURI(originalURLStruct.OriginalURL); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.Log.Error("Error:",
				zap.String("reason", "The request body isn't a url"))
			return
		}
	}

	//Вызов метода SetListData для полученных данных
	resp, err := h.Storage.SetListData(r.Context(), req, h.prefix)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
	//Формирование ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err = enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
}