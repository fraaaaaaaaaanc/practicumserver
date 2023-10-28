// Package storage defines the StorageMock interface, which outlines the methods for interacting with the data storage.
package storage

import (
	"context"
	"practicumserver/internal/models"
)

type StorageMock interface {
	// SetData accepts an original URL and stores it in the storage.
	// It returns the shortened URL for the provided original URL and an error if any issue occurs.
	SetData(ctx context.Context, link string) (string, error)

	// GetData receives a shortened URL and retrieves the associated original URL from the storage.
	// It returns the original URL and an error if the shortened URL is found.
	// If the URL is not found, it returns an empty string.
	GetData(ctx context.Context, shortLink string) (string, error)

	// SetListData takes a slice of original URLs and stores each URL in the storage.
	// It returns a slice of models.ResponseAPIBatch, which contains shortened URLs for the original URLs, and an error.
	SetListData(ctx context.Context, reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error)

	// GetListData retrieves a list of all URLs submitted by the user.
	// It returns a slice of models.ResponseAPIUserUrls and an error.
	GetListData(ctx context.Context, prefix string) ([]models.ResponseAPIUserUrls, error)

	// CheckUserID checks if a generated UserID is unique within the storage.
	// It returns true if the UserID is unique, false if it exists, and an error if any issue occurs.
	CheckUserID(ctx context.Context, userID string) (bool, error)

	// UpdateDeletedFlag modifies the deletion flag of URLs based on the user and shortLink lists.
	// It does not return any data but may return an error if there's an issue.
	UpdateDeletedFlag(ctx context.Context, userIDList, shortLinkList []string) error
}
