package router

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/compress"
	"practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/storage"
)

func Router(log *zap.Logger, prefix string,
	dbStorageAdress, fileStoragePath string) (chi.Router, error) {
	//Создание объекта storage реализующего интерфейсный тип storage.StorageMock
	strg, err := storage.NewStorage(log, dbStorageAdress, fileStoragePath)
	if err != nil {
		return nil, err
	}
	//Создание объекта handlers
	hndlrs := handlers.NewHandlers(strg, log, prefix)

	//Создание объекта *Mux
	r := chi.NewRouter()

	r.Use(logger.MiddlewareLogHandleFunc(log), compress.MiddlewareGzipHandleFunc(log))
	r.Get("/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		//Хендлер для получения оригинального URL по сокращенному
		hndlrs.GetRequest(w, r)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		//Хендлер для проверки подключения к бд
		hndlrs.GetRequestPing(w, r)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		//Хендлер для получения сокращенного URL
		hndlrs.PostRequest(w, r)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		//Хендлер для получения сокращенного URL, этот хендлер умеет принимать
		//и обрабатывать запросы в формате JSON
		hndlrs.PostRequestAPIShorten(w, r)
	})
	r.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		//Хендлер для получения сокращенного URL, этот хендлер умеет принимать
		//и обрабатывать запросы в формате JSON пачками
		hndlrs.PostRequestAPIShortenBatch(w, r)
	})

	return r, nil
}
