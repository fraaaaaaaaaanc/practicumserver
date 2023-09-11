package handlers

import (
	"io"
	"net/http"
	"practicumserver/internal/config"
	"practicumserver/internal/db"
	utils2 "practicumserver/internal/utils"
)

var encodigs []string = []string{"charset=utf-8", "charset=iso-8859-1", "charset=windows-1251", "charset=us-ascii"}

// Обработчик Post запроса
func PostRequest(w http.ResponseWriter, r *http.Request, flags *config.Flags) {
	contentType := r.Header.Get("Content-Type")
	if !utils2.ValidContentType(contentType) || r.URL.String() != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	link, err := io.ReadAll(r.Body)
	if err != nil || string(link) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	srtLink := utils2.LinkShortening()
	avlblSrtLink, err := db.SetDB(string(link), srtLink)
	if err != nil {
		srtLink = avlblSrtLink
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(flags.ShortLink + "/" + srtLink))
}

// Обработчик Get запроса
func GetRequest(w http.ResponseWriter, r *http.Request) {
	link := r.URL.String()[1:]
	baseLink, err := db.GetDB(link)
	if link == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", baseLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
