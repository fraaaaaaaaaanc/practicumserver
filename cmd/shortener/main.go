package main

import (
	"fmt"
	"net/http"
	"practicumserver/cmd/shortener/config"
	"practicumserver/cmd/shortener/router"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flags := config.ParseConfFlugs()

	fmt.Println("Running server on", flags.String())
	return http.ListenAndServe(flags.String(), router.Router(flags))
}
