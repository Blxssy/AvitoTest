package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(logLevel string) *zap.Logger {
	config := zap.NewProductionEncoderConfig()

	config.MessageKey = "message"
	config.TimeKey = "timestamp"
	config.StacktraceKey = ""
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewJSONEncoder(config)

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		return zap.New(zapcore.NewCore(consoleEncoder, os.Stdout, zap.DebugLevel))
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	return zap.New(core, zap.AddCaller())
}
