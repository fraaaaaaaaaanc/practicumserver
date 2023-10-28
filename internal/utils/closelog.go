package utils

import (
	"go.uber.org/zap"
	"practicumserver/internal/config"
	"practicumserver/internal/logger"
)

// Closelog is a function to close the log resources.
// It takes a logger and configuration flags as input and performs the following tasks:
// 1. Synchronizes the logger to ensure all log entries are flushed.
// 2. If the 'FileLog' flag is set in the configuration, it attempts to close the log file and logs an error if it fails.
func Closelog(log *logger.ZapLogger, flags *config.Flags) {
	log.Logger.Sync()
	if flags.FileLog {
		if err := log.File.Close(); err != nil {
			log.Error("Failed to close log file", zap.Error(err))
		}
	}
}
