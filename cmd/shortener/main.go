package main

import (
	"net/http"
	"practicumserver/cmd/shortener/router"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, router.Router())
}
