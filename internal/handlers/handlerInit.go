package handlers

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"practicumserver/internal/storage"
)

type Handlers struct {
	Storage storage.StorageMock
	Log     *zap.Logger
	prefix  string
}

func NewHandlers(strg storage.StorageMock, log *zap.Logger, prefix string) *Handlers {
	return &Handlers{
		Storage: strg,
		Log:     log,
		prefix:  prefix,
	}
}
