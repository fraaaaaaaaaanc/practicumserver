package models

type (
	RequestApiShorten struct {
		LongURL string `json:"url"`
	}

	ResponseApiShorten struct {
		ShortURL string `json:"result"`
	}
)

type (
	RequestApiBatch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalUrl   string `json:"original_url"`
	}

	ResponseApiBatch struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)
