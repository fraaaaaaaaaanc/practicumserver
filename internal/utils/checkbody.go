// Package utils provides a collection of utility functions and tools
// used across the LinksShortener Server application. These utilities cover a wide
// range of common tasks, such as handling request bodies, generating random
// // shortened links, closing logs, and more. This package serves to improve
// code modularity and maintainability by centralizing common utility functions.
package utils

import (
	"bytes"
	"io"
)

// IsRequestBodyEmpty checks if the request body is empty.
// It reads the body and returns true if it's empty, false otherwise.
// If there's an error while reading the body, it returns an error.
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
