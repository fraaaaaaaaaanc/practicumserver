package main

import (
	"fmt"
	"net/http"
	"practicumserver/cmd/shortener/config"
	"practicumserver/cmd/shortener/router"
)

func main() {
	flags := config.ParseFlags()

	if err := run(flags); err != nil {
		panic(err)
	}
}

func run(flags *config.Flags) error {
	fmt.Println("Running server on", flags.HostFlag)
	return http.ListenAndServe(flags.HostFlag, router.Router(flags))
}
