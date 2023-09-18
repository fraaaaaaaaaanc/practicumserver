package main

import (
	"fmt"
	"net/http"
	"practicumserver/internal/config"
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
	var storage storage2.Storage
	storage.ShortUrls = map[string]string{
		"test": "http://test",
	}
	storage.ShortBoolUrls = map[string]bool{
		"test": true,
	}

	fmt.Println("Running server on", flags.String())
	return http.ListenAndServe(flags.String(), router.Router(flags, storage))
}
