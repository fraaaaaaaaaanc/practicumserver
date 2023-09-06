package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"practicumserver/cmd/shortener/handlers"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.PostRequest).Methods("POST")
	router.HandleFunc("/{id:[a-zA-Z0-9]+}", handlers.GetRequest).Methods("GET")

	return http.ListenAndServe(`:8080`, router)
}
