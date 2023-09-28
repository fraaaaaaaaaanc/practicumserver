package models

type (
	Request struct {
		LongUrl string `json:"url"`
	}

	Response struct {
		ShortUrl string `json:"result"`
	}
)
