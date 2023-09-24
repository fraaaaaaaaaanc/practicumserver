package main

import (
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/config"
	"practicumserver/internal/logger"
	"practicumserver/internal/router"
	storage2 "practicumserver/internal/storage"
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
	log := logger.NewZapLogger(flags)
	storage := storage2.NewStorage()

	defer utils.Closelog(log, flags)

	log.Info("Server start", zap.String("Running server on", flags.String()))
	return http.ListenAndServe(flags.String(), router.Router(flags, storage, log.Logger))
}
