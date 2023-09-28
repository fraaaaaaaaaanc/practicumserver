package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
	"practicumserver/internal/utils"
)

var encodigs []string = []string{"charset=utf-8", "charset=iso-8859-1", "charset=windows-1251", "charset=us-ascii"}

type Handlers struct {
}

// Обработчик Post запроса
func (h *Handlers) PostRequest(w http.ResponseWriter, r *http.Request, storage *storage.Storage, flags string) {
	contentType := r.Header.Get("Content-Type")
	if !utils.ValidContentType(contentType, "text/plain") || r.URL.String() != "/" {
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
func (h *Handlers) GetRequest(w http.ResponseWriter, r *http.Request, storage *storage.Storage) {
	shortLink := r.URL.String()[1:]
	baseLink, boolRes := storage.GetData(shortLink)
	if shortLink == "" || boolRes {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", baseLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) PostRequestAPIShorten(w http.ResponseWriter, r *http.Request, strg *storage.Storage, flag string) {
	contentType := r.Header.Get("Content-Type")
	if !utils.ValidContentType(contentType, "application/json") ||
		r.URL.String() != "/api/shorten" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if req.LongURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink := strg.GetNewShortLink(req.LongURL)
	strg.SetData(req.LongURL, shortLink)

	resp := models.Response{
		ShortURL: flag + "/" + shortLink,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
