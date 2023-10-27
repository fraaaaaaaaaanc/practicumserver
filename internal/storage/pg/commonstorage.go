package pgstorage

import (
	"go.uber.org/zap"
	"sync"
)

// Структура с ощими элементами для каждого storage
type StorageParam struct {
	Log *zap.Logger
	Sm  *sync.Mutex
}

type GetResponse struct {
	originalURL string
	deletedFlag bool
}
