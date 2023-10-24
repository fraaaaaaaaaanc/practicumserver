package utils

import (
	"bytes"
	"encoding/json"
)

// Функция которая вызывалась в MiddlewareGzipHandleFunc, для определения формата полученных данных,
// так как запрос в тестах прилетает с contentType-ом application/x-gzip, в последствии оказалось,
// что так делать не нужно
func IsJSON(data []byte) bool {
	var js interface{}
	if err := json.Unmarshal(data, &js); err != nil {
		return false
	}
	return true
}

func IsText(data []byte) bool {
	controlChars := "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x0B\x0C\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"
	for _, char := range controlChars {
		if bytes.IndexByte(data, byte(char)) != -1 {
			return false
		}
	}
	return true
}

func DetermineContentType(data []byte) string {
	if IsJSON(data) {
		return "application/json"
	}
	if IsText(data) {
		return "text/plain"
	}
	return ""
}
