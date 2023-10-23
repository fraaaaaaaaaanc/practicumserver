package models

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

type ResponseApiUserUrls struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
