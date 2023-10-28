// Package models defines data structures and error constants used throughout the application.
// It includes structures for request and response formats, error definitions, and custom data types
// for handling data related to user URLs and requests. This package provides a central location
// for defining and managing data structures used in various parts of the application.
package models

import "net/http"

type ContextKey string

// ContextKey is a custom type representing a key for storing and retrieving values in the context.
var UserIDKey ContextKey = "userID"

// DeleteURLList is a slice of strings representing a list of short URLs to be deleted.
type DeleteURLList []string

type (
	// RequestAPIShorten is a struct representing the request format for a POST request to "/api/shorten".
	RequestAPIShorten struct {
		OriginalURL string `json:"url"` // OriginalURL is the URL to be shortened.
	}
	// ResponseAPIShorten is a struct representing the response format for a POST request to "/api/shorten".
	ResponseAPIShorten struct {
		ShortLink string `json:"result"` // ShortLink is the shortened URL result.
	}
)

type (
	// RequestAPIBatch is a struct representing the request format for a batch POST request to "/api/shorten/batch".
	RequestAPIBatch struct {
		CorrelationID string `json:"correlation_id"` // CorrelationID is an identifier for the batch request.
		OriginalURL   string `json:"original_url"`   // OriginalURL is the URL to be shortened.
	}
	// ResponseAPIBatch is a struct representing the response format for a batch POST request to "/api/shorten/batch".
	ResponseAPIBatch struct {
		CorrelationID string `json:"correlation_id"` // CorrelationID is the identifier for the batch request.
		ShortURL      string `json:"short_url"`      // ShortURL is the shortened URL result.
	}
)

// ResponseAPIUserUrls is a struct representing the response format for a GET request to retrieve user URLs.
type ResponseAPIUserUrls struct {
	ShortURL    string `json:"short_url"`    // ShortURL is the shortened URL.
	OriginalURL string `json:"original_url"` // OriginalURL is the original URL.
}

// FileData is a struct representing data stored in a file with user information.
type FileData struct {
	UserID      string `json:"user_id"`      // UserID is the user identifier.
	ShortURL    string `json:"short_url"`    // ShortURL is the shortened URL.
	OriginalURL string `json:"original_url"` // OriginalURL is the original URL.
	DeletedFlag bool   `json:"deleted_flag"` // DeletedFlag indicates if the data has been deleted.
}

// DeleteURL is a struct representing the data required for deleting a user's URL.
type DeleteURL struct {
	UserID   string // UserID is the user identifier.
	ShortURL string // ShortURL is the shortened URL.
}

// HandlerFuncAdapter type for type identification func(http.ResponseWriter, *http.Request)
type HandlerFuncAdapter func(http.ResponseWriter, *http.Request)

func (h HandlerFuncAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}
