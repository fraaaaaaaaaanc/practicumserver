package models

type (
	RequestAPIShorten struct {
		LongURL string `json:"url"`
	}

	ResponseAPIShorten struct {
		ShortURL string `json:"result"`
	}
)

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
