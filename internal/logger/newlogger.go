package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"practicumserver/internal/config"
)

type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
}

type ZapLogger struct {
	Logger *zap.Logger
	File   *os.File
}

func NewZapLogger(flags *config.Flags) *ZapLogger {
	var cores []zapcore.Core
	var file *os.File

	consoleConfig := zap.NewDevelopmentConfig()
	consoleConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(consoleConfig.EncoderConfig),
		zapcore.Lock(os.Stdout),
		zapcore.InfoLevel))

	if flags.FileLog {
		file, err := os.OpenFile("filelogger.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}

		writeSyncer := zapcore.AddSync(file)

		fileConfig := zap.NewProductionConfig()
		fmt.Println(file.Name())
		//fileConfig.OutputPaths = []string{file.Name()}
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(fileConfig.EncoderConfig),
			writeSyncer,
			zapcore.InfoLevel))
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller())
	return &ZapLogger{Logger: logger, File: file}
}

func (z *ZapLogger) Info(msg string, fileds ...zapcore.Field) {
	z.Logger.Info(msg, fileds...)
}

func (z *ZapLogger) Error(msg string, fileds ...zapcore.Field) {
	z.Logger.Error(msg, fileds...)
}
