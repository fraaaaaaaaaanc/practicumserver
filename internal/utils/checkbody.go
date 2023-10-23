package utils

import (
	"bytes"
	"io"
)

// Функция проверяющая тело запроса на пустоту
func IsRequestBodyEmpty(body io.Reader) (bool, error) {
	newBody, err := io.ReadAll(body)
	if err != nil {
		return false, err
	}
	body = io.NopCloser(bytes.NewReader(newBody))
	if newBody == nil {
		return true, nil
	}
	return false, nil
}
