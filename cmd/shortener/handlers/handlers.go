package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"practicumserver/cmd/shortener/db"
	"practicumserver/cmd/shortener/utils"
)

var encodigs []string = []string{"charset=utf-8", "charset=iso-8859-1", "charset=windows-1251", "charset=us-ascii"}

// Обработчик Post запроса
func PostRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !utils.ValidContentType(contentType) || r.URL.String() != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	link, err := io.ReadAll(r.Body)
	if err != nil || string(link) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	srtLink := utils.LinkShortening()
	avlblSrtLink, err := db.SetDB(string(link), srtLink)
	if err != nil {
		srtLink = avlblSrtLink
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	domain := fmt.Sprintf("http://%s/", r.Host)
	_, _ = w.Write([]byte(domain + srtLink))
}

// Обработчик Get запроса
func GetRequest(w http.ResponseWriter, r *http.Request) {
	link := r.URL.String()[1:]
	baseLink, err := db.GetDB(link)
	fmt.Println(link, baseLink, err)
	if link == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", baseLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", GetRequest)
	r.Post("/", PostRequest)

	return r
}
