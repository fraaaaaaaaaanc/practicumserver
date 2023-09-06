package handlers

import (
	"io"
	"net/http"
	"practicumserver/cmd/shortener/db"
	"practicumserver/cmd/shortener/utils"
)

// Локальный адресс
var LocalURL string = "http://localhost:8080/"

// Обработчик Post запроса
func PostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" || r.URL.String() != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil || string(body) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	id := utils.LinkShortening()
	if _, ok := db.ShortUrls[id]; !ok {
		db.ShortUrls[id] = string(body)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(LocalURL + id))
}

// Обработчик Get запроса
func GetRequest(w http.ResponseWriter, r *http.Request) {
	if _, ok := db.ShortUrls[r.URL.String()[1:]]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", db.ShortUrls[r.URL.String()[1:]])
	w.WriteHeader(http.StatusTemporaryRedirect)
}
