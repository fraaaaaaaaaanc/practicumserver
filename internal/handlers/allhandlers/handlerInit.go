// Package handlers provides HTTP request handlers and middleware for the LinksShortener application.
// These handlers are responsible for processing incoming HTTP requests, performing necessary actions,
// and generating appropriate HTTP responses. They interact with the application's business logic and
// data storage to implement various functionalities such as URL shortening, retrieving user URLs,
// and handling batch URL shortening requests.
// The package contains handlers for different routes and HTTP methods, each with a specific purpose.
package handlers

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
)

// Handlers is a structure that encapsulates the application's request handlers and related components.
type Handlers struct {
	Storage storage.StorageMock
	Log     *zap.Logger
	DelCn   chan *models.DeleteURL
	prefix  string
}

// New Handlers initializes and returns a Handlers object. It takes a storage object that implements the
// storage.StorageMock interface, a logger object for logging,
// and a prefix string obtained from the -b flag for generating responses in POST handlers.
func NewHandlers(strg storage.StorageMock, log *zap.Logger, prefix string) *Handlers {
	return &Handlers{
		Storage: strg,                              // Storage is the storage component implementing the storage.StorageMock interface.
		Log:     log,                               // Log is the logger used for logging within the request handlers.
		DelCn:   make(chan *models.DeleteURL, 128), // DelCn is a channel for handling delete URL requests.
		prefix:  prefix,
	}
}
