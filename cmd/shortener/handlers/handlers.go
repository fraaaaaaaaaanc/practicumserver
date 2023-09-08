package handlers

import (
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"practicumserver/cmd/shortener/db"
	"practicumserver/cmd/shortener/utils"
)

var encodigs []string = []string{"charset=utf-8", "charset=iso-8859-1", "charset=windows-1251", "charset=us-ascii"}

// Локальный адресс
var LocalURL string = "http://localhost:8080/"

// Обработчик Post запроса
func PostRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if utils.ValidContentType(contentType) || r.URL.String() != "/" {
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
	if _, ok := db.ShortUrls[string(body)]; !ok {
		db.ShortUrls[string(body)] = id
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(LocalURL + id))
}

// Обработчик Get запроса
func GetRequest(w http.ResponseWriter, r *http.Request) {
	for k, v := range db.ShortUrls {
		if v == r.URL.String()[1:] {
			w.Header().Set("Location", k)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", GetRequest)
	r.Post("/", PostRequest)

	return r
}
