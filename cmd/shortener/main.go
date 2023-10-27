package main

import (
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/app"
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
	appStrct, err := app.NewApp()
	if err != nil {
		return err
	}
	//Закрытие логов
	defer utils.Closelog(appStrct.Log, appStrct.Flags)

	defer func() {
		if DBstrg, ok := appStrct.Strg.(*pgstorage.DBStorage); ok {
			DBstrg.DB.Close()
		}
	}()

	//Запуск сервера по адресу переданному через flags.String()
	appStrct.Log.Info("Server start", zap.String("Running server on:", appStrct.Flags.String()))
	err = http.ListenAndServe(appStrct.Flags.String(), appStrct.Rtr)
	appStrct.Log.Error("Error:", zap.Error(err))
	return err
}
