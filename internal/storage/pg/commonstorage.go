// Package pgstorage contains implementations of various storage backends.
package pgstorage

import (
	"go.uber.org/zap"
	"sync"
)

// StorageParam is a structure with common elements for each storage.
type StorageParam struct {
	Log *zap.Logger // Log is a logger used for logging operations and errors.
	Sm  *sync.Mutex // Sm is a mutex for synchronizing access to shared resources in storage.
}

// GetResponse is a structure representing the response for the GetData method.
type GetResponse struct {
	originalURL string
	deletedFlag bool
}
