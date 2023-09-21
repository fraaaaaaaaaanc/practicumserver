package main

import (
	"fmt"
	"net/http"
	"practicumserver/internal/config"
	"practicumserver/internal/logger"
	"practicumserver/internal/router"
	storage2 "practicumserver/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flags := config.ParseConfFlugs()
	if err := logger.Initialize(flags.LogLevel); err != nil {
		return err
	}

	storage := *storage2.NewStorage()

	fmt.Println("Running server on", flags.String())
	return http.ListenAndServe(flags.String(), router.Router(flags, storage))
}
