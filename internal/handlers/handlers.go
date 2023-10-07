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
	Storage     *storage.Storage
	shortLink   string
	fileStorage string
}

func NewHandlers(shortLink, fileStorage string) *Handlers {
	strg := storage.NewStorage()

	return &Handlers{
		Storage:     strg,
		shortLink:   shortLink,
		fileStorage: fileStorage,
	}
}

// Обработчик Post запроса
func (h *Handlers) PostRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if (!utils.ValidContentType(contentType, "text/plain") && contentType != "application/x-gzip") ||
		r.URL.String() != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	link, err := io.ReadAll(r.Body)
	if err != nil || string(link) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newShortLink := h.Storage.GetNewShortLink(string(link), h.fileStorage)
	h.Storage.SetData(string(link), newShortLink)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(h.shortLink + "/" + newShortLink))
}

// Обработчик Get запроса
func (h *Handlers) GetRequest(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.String()[1:]
	baseLink, boolRes := h.Storage.GetData(shortLink)
	if shortLink == "" || boolRes {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", baseLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) PostRequestAPIShorten(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if (!utils.ValidContentType(contentType, "application/json") && contentType != "application/x-gzip") ||
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

	newShortLink := h.Storage.GetNewShortLink(req.LongURL, h.fileStorage)
	h.Storage.SetData(req.LongURL, newShortLink)

	resp := models.Response{
		ShortURL: h.shortLink + "/" + newShortLink,
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
