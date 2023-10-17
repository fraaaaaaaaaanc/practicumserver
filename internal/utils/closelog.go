package utils

import (
	"go.uber.org/zap"
	"practicumserver/internal/config"
	"practicumserver/internal/logger"
)

// Функция закрывающая логи
func Closelog(log *logger.ZapLogger, flags *config.Flags) {
	log.Logger.Sync()
	if flags.FileLog {
		if err := log.File.Close(); err != nil {
			log.Error("Failed to close log file", zap.Error(err))
		}
	}
}
