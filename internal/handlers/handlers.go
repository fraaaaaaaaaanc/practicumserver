package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
	"practicumserver/internal/utils"
	"time"
)

type Handlers struct {
	Storage     storage.StorageMock
	Log         *zap.Logger
	shortLink   string
	fileStorage string
	dbAdress    string
}

func NewHandlers(strg storage.StorageMock, log *zap.Logger, shortLink, dbAdress, fileStorage string) *Handlers {
	return &Handlers{
		Storage:     strg,
		Log:         log,
		shortLink:   shortLink,
		fileStorage: fileStorage,
		dbAdress:    dbAdress,
	}
}

// Обработчик Post запроса
func (h *Handlers) PostRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if (!utils.ValidContentType(contentType, "text/plain") && contentType != "application/x-gzip") ||
		r.URL.String() != "/" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL or Content-Type"))
		return
	}
	link, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	if _, err := url.Parse(string(link)); err != nil || string(link) == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "The request body isn't a url"))
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			//w.WriteHeader(http.StatusInternalServerError)
			//h.Log.Error("Error:", zap.Error(err))
			return
		}
	}()

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
	if boolRes {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error: This abbreviated link wasn't found",
			zap.String("shortLink", shortLink))
		return
	}
	w.Header().Set("Location", baseLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) GerRequestPing(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", h.dbAdress)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
	ctx, cansel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cansel()
	if err = db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:", zap.Error(err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) PostRequestAPIShorten(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if (!utils.ValidContentType(contentType, "application/json") && contentType != "application/x-gzip") ||
		r.URL.String() != "/api/shorten" {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "Invalid URL or Content-Type"))
		return
	}

	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error:", zap.Error(err))
		return
	}

	if _, err := url.Parse(req.LongURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error:",
			zap.String("reason", "The request body isn't a url"))
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
		h.Log.Error("Error:", zap.Error(err))
		return
	}
}
