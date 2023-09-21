package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	log = zl.Sugar()

	return nil
}
