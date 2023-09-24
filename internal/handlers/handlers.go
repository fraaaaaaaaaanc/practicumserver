package handlers

import (
	"io"
	"net/http"
	storage2 "practicumserver/internal/storage"
	utils2 "practicumserver/internal/utils"
)

var encodigs []string = []string{"charset=utf-8", "charset=iso-8859-1", "charset=windows-1251", "charset=us-ascii"}

type Handlers struct {
}

// Обработчик Post запроса
func (h *Handlers) PostRequest(w http.ResponseWriter, r *http.Request, storage *storage2.Storage, flags string) {
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

	shortLink := storage.GetNewShortLink(string(link))
	storage.SetData(string(link), shortLink)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(flags + "/" + shortLink))
}

// Обработчик Get запроса
func (h *Handlers) GetRequest(w http.ResponseWriter, r *http.Request, storage *storage2.Storage) {
	shortLink := r.URL.String()[1:]
	baseLink, boolRes := storage.GetData(shortLink)
	if shortLink == "" || boolRes {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", baseLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
