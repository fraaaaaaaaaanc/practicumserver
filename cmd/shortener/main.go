package main

import (
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/config"
	handlers "practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/router"
	"practicumserver/internal/storage/pg"
	"practicumserver/internal/utils"
)

func main() {
	//Начало работы программы
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	//Парсинг флагов и создание переменной логирования
	flags := config.ParseConfFlugs()
	log, err := logger.NewZapLogger(flags.FileLog)
	if err != nil {
		return err
	}
	//Закрытие логов
	defer utils.Closelog(log, flags)
	//Создание объекта storage реализующего интерфейсный тип storage.StorageMock
	strg, err := storage.NewStorage(log.Logger, flags.DBStorageAdress, flags.FileStoragePath)
	if err != nil {
		return err
	}
	defer func() {
		if DBstrg, ok := strg.(*storage.DBStorage); ok {
			DBstrg.DB.Close()
		}
	}()
	//Создание объекта handlers
	hndlrs := handlers.NewHandlers(strg, log.Logger, flags.Prefix)
	//Создание объекта роутера для передачи в http.ListenAndServe
	rtr, err := router.Router(hndlrs, log.Logger)
	if err != nil {
		return err
	}

	//Запуск сервера по адресу переданному через flags.String()
	log.Info("Server start", zap.String("Running server on:", flags.String()))
	err = http.ListenAndServe(flags.String(), rtr)
	log.Error("Error:", zap.Error(err))
	return err
}
