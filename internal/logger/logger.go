package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the global logger instance
var Logger *zap.Logger

// Init initializes the logger based on configuration
func Init(level, format, output string) error {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Set log level
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		logLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(logLevel)

	// Set output
	if output != "stdout" && output != "stderr" {
		config.OutputPaths = []string{output}
		config.ErrorOutputPaths = []string{output}
	}

	// Build logger
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	// Replace global logger
	zap.ReplaceGlobals(Logger)

	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if Logger == nil {
		// Fallback to default logger if not initialized
		Logger, _ = zap.NewProduction()
	}
	return Logger
}

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// WithRequestID creates a logger with request ID
func WithRequestID(requestID string) *zap.Logger {
	return Logger.With(zap.String("request_id", requestID))
}
