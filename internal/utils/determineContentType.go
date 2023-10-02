package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func IsText(data []byte) bool {
	controlChars := "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x0B\x0C\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"
	for _, char := range controlChars {
		if bytes.IndexByte(data, byte(char)) != -1 {
			return false
		}
	}
	return true
}

func isJSON(data []byte) bool {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return true
	}
	return false
}

func DetermineContentType(req *http.Request) (string, error) {
	copyData, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	req.Body = io.NopCloser(bytes.NewReader(copyData))
	if isJSON(copyData) {
		return "application/json", nil
	}
	if IsText(copyData) {
		return "text/plain", nil
	}
	return "", err
}
