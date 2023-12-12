package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new logger with the desired structure
func New(level string) (*zap.Logger, error) {
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("could not parse level: %w", err)
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.Lock(os.Stdout),
		l,
	))

	return logger, nil
}
