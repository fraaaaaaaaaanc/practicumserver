package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Logger is the interface that a logger should implement.
type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
}

// ZapLogger is a logger structure that includes a logger and a file to write logs.
type ZapLogger struct {
	Logger *zap.Logger // Logger a pointer to a Zap logger that handles the logging functionality.
	File   *os.File    // File a pointer to an os.File that, if not nil, represents a file where logs are written in addition to the console.
}

// NewZapLogger initializes a logger. It accepts a boolean parameter, FileLog, which is obtained when parsing flags.
// If FileLog is true, the logger instance is created to write logs to both the console and a file.
// Otherwise, logs are written only to the console.
func NewZapLogger(FileLog bool) (*ZapLogger, error) {
	var cores []zapcore.Core
	var file *os.File

	consoleConfig := zap.NewDevelopmentConfig()
	consoleConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(consoleConfig.EncoderConfig),
		zapcore.Lock(os.Stdout),
		zapcore.InfoLevel))

	if FileLog {
		file, err := os.OpenFile("filelogger.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}

		writeSyncer := zapcore.AddSync(file)

		fileConfig := zap.NewProductionConfig()
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(fileConfig.EncoderConfig),
			writeSyncer,
			zapcore.InfoLevel))
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller())
	return &ZapLogger{Logger: logger, File: file}, nil
}

// Info is a method to log informational messages.
func (z *ZapLogger) Info(msg string, fields ...zapcore.Field) {
	z.Logger.Info(msg, fields...)
}

// Error is a method to log error messages.
func (z *ZapLogger) Error(msg string, fields ...zapcore.Field) {
	z.Logger.Error(msg, fields...)
}
