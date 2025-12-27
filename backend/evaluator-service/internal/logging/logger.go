package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct{ *zap.Logger }

func New() *Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, _ := cfg.Build()
	return &Logger{l}
}

func (l *Logger) Info(msg string, fields ...zap.Field)  { l.Logger.Info(msg, fields...) }
func (l *Logger) Warn(msg string, fields ...zap.Field)  { l.Logger.Warn(msg, fields...) }
func (l *Logger) Error(msg string, fields ...zap.Field) { l.Logger.Error(msg, fields...) }
func (l *Logger) Fatal(msg string, fields ...zap.Field) { l.Logger.Fatal(msg, fields...) }

func KV(k string, v interface{}) zap.Field { return zap.Any(k, v) }
func Err(err error) zap.Field              { return zap.Error(err) }
