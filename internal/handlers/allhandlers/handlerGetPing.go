package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/storage/pg"
)

// Хендлер принимающий GET для проверки соединения с базой данных
func (h *Handlers) GetRequestPing(w http.ResponseWriter, r *http.Request) {
	//Проверяемм пренадлежит ли хранилище инициализированное в hndlrs типу *storage.DBStorage
	if ds, ok := h.Storage.(*storage.DBStorage); ok {
		//Если да то выщываем функцию PingDB для проверки соединения
		if err := ds.PingDB(r.Context()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.Log.Error("Error:", zap.Error(err))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusBadRequest)
}
