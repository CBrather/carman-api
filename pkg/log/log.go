package log

import (
	"time"

	"github.com/go-logr/zapr"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Initializes the global logger and other lib-specific loggers with custom field formats and the given log level
func Initialize(level string) error {
	logger, err := GetLoggerWithLevel(level)
	if err != nil {
		zap.L().Fatal("Failed to initialize the global logger", zap.Error(err))
	}

	zap.ReplaceGlobals(logger)

	otelLogger := zapr.NewLogger(logger)
	otel.SetLogger(otelLogger)

	return nil
}

// Returns a new logger with custom field formats and the given log level (Only to be used when the global one cannot)
func GetLoggerWithLevel(level string) (*zap.Logger, error) {
	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		parsedLevel = zap.InfoLevel
	}

	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(parsedLevel)
	setEncoderConfig(&config)

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func setEncoderConfig(cfg *zap.Config) {
	cfg.EncoderConfig.EncodeTime = timeEncoder
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.StacktraceKey = "callstack"
	cfg.EncoderConfig.TimeKey = "timestamp"
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format(time.RFC3339Nano))
}
