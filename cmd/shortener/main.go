package main

import (
	"github.com/gorilla/mux"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var ShortUrls = map[string]string{
	"test": "http://test",
}
var LocalURL string = "http://localhost:8080/"

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	router := mux.NewRouter()
	router.HandleFunc("/", PostRequest).Methods("POST")
	router.HandleFunc("/{id:[a-zA-Z0-9]+}", GetRequest).Methods("GET")

	return http.ListenAndServe(`:8080`, router)
}

func PostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" || r.URL.String() != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil || string(body) == "" {
		//http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	//rand.Seed(time.Now().UnixNano())
	id := Base62Encode(rand.Uint64())
	if _, ok := ShortUrls[id]; !ok {
		ShortUrls[id] = string(body)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(LocalURL + id))
}

func GetRequest(w http.ResponseWriter, r *http.Request) {
	if _, ok := ShortUrls[r.URL.String()[1:]]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		//r.Header.Get("Content-Type") != "text/plain" ||
		return
	}
	//w.Header().Set("Content-Type", "text/plain")
	//for key, values := range w.Header() {
	//	for _, value := range values {
	//		fmt.Println(w, "%s: %s\n", key, value)
	//	}
	//}
	//_, _ = w.Write([]byte{})
	w.Header().Set("Location", ShortUrls[r.URL.String()[1:]])
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func Base62Encode(number uint64) string {
	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}

	return encodedBuilder.String()
}

//func webhook(w http.ResponseWriter, r *http.Request) {
//	switch r.Method {
//	case http.MethodGet:
//		GetRequest(w, r)
//	case http.MethodPost:
//		if err := r.ParseForm(); err != nil {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//		PostRequest(w, r)
//	}
//	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "text/plain" {
//		//w.Header().Set("Content-Type", "application/json")
//		//// пока установим ответ-заглушку, без проверки ошибок
//		//_, _ = w.Write([]byte(`
//		// {
//		//   "response": {
//		//     "text": "Извините, я пока ничего не умею"
//		//   },
//		//   "version": "1.0"
//		// }`))
//		_, err := io.ReadAll(r.Body)
//		if err != nil {
//			http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
//			return
//		}
//		defer r.Body.Close()
//
//		w.Header().Set("Content-Type", "text/plain")
//		_, _ = w.Write([]byte(r.Header.Get("Content-Length")))
//	} else {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//}
