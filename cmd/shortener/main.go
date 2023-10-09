package main

import (
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/config"
	"practicumserver/internal/logger"
	"practicumserver/internal/router"
	"practicumserver/internal/utils"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	//creating an instance of flags, storage, logs
	flags := config.ParseConfFlugs()
	log, err := logger.NewZapLogger(flags.FileLog)
	if err != nil {
		return err
	}

	defer utils.Closelog(log, flags)

	rtr, err := router.Router(log.Logger, flags.Prefix, flags.FileStorage, flags.DBAdress)
	if err != nil {
		return err
	}

	log.Info("Server start", zap.String("Running server on:", flags.String()))
	err = http.ListenAndServe(flags.String(), rtr)
	log.Error("Error:", zap.Error(err))
	return err
}
