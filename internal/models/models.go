package models

type ContextKey string

var UserIDKey ContextKey = "userID"

type DeleteURLList []string

// Стркуры реализующие формы запросов и ответов при POST запросах к адрессу /api/shorten
type (
	RequestAPIShorten struct {
		OriginalURL string `json:"url"`
	}

	ResponseAPIShorten struct {
		ShortLink string `json:"result"`
	}
)

// Стркуры реализующие формы запросов и ответов при POST запросах к адрессу /api/shorten/batch
type (
	RequestAPIBatch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	ResponseAPIBatch struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)

type ResponseAPIUserUrls struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileData struct {
	UserID      string `json:"user_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	DeletedFlag bool   `json:"deleted_flag"`
}

type DeleteURL struct {
	UserID   string
	ShortURL string
}
