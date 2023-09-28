package models

type (
	Request struct {
		LongURL string `json:"url"`
	}

	Response struct {
		ShortURL string `json:"result"`
	}
)
