package handlers

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
)

// Струкутра объекта Handlers
type Handlers struct {
	Storage storage.StorageMock
	Log     *zap.Logger
	DelCn   chan *models.DeleteURL
	prefix  string
}

// Инициализатор объекта Handlers, принимающий объект хранилища реализующий интерфейс storage.StorageMock,
// объект log для логирования и строку prefix полученная из флага -b для формирования ответа в POST хендлерах
func NewHandlers(strg storage.StorageMock, log *zap.Logger, prefix string) *Handlers {
	return &Handlers{
		Storage: strg,
		Log:     log,
		DelCn:   make(chan *models.DeleteURL, 128),
		prefix:  prefix,
	}
}
