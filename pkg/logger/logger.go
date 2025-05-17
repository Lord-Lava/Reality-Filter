package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Global logger instance
	globalLogger *zap.Logger
	once         sync.Once
)

// Config holds the configuration for the logger
type Config struct {
	// LogLevel is the minimum enabled logging level
	LogLevel string `yaml:"level" json:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally
	Development bool `yaml:"development" json:"development"`
	// Encoding sets the logger's encoding. Valid values are "json" and "console"
	Encoding string `yaml:"encoding" json:"encoding"`
	// OutputPaths is a list of URLs or file paths to write logging output to
	OutputPaths []string `yaml:"output_paths" json:"output_paths"`
	// ErrorOutputPaths is a list of URLs to write internal logger errors to
	ErrorOutputPaths []string `yaml:"error_output_paths" json:"error_output_paths"`
}

// DefaultConfig returns a default configuration for the logger
func DefaultConfig() Config {
	return Config{
		LogLevel:         "info",
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// parseLevel parses a level string into a zapcore.Level
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Init initializes the global logger with the given configuration
func Init(cfg Config) error {
	var err error
	once.Do(func() {
		// Create basic encoder config
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// Create zap configuration
		zapConfig := zap.Config{
			Level:            zap.NewAtomicLevelAt(parseLevel(cfg.LogLevel)),
			Development:      cfg.Development,
			Encoding:         cfg.Encoding,
			EncoderConfig:    encoderConfig,
			OutputPaths:      cfg.OutputPaths,
			ErrorOutputPaths: cfg.ErrorOutputPaths,
		}

		// Build the logger
		globalLogger, err = zapConfig.Build(
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
		if err != nil {
			fmt.Printf("Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
	})

	return err
}

// Sync flushes any buffered log entries
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// With creates a child logger with the given fields
func With(fields ...zap.Field) *zap.Logger {
	if globalLogger == nil {
		Init(DefaultConfig())
	}
	return globalLogger.With(fields...)
}

// Debug logs a message at DebugLevel
func Debug(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		Init(DefaultConfig())
	}
	globalLogger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel
func Info(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		Init(DefaultConfig())
	}
	globalLogger.Info(msg, fields...)
}

// Warn logs a message at WarnLevel
func Warn(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		Init(DefaultConfig())
	}
	globalLogger.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel
func Error(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		Init(DefaultConfig())
	}
	globalLogger.Error(msg, fields...)
}

// Fatal logs a message at FatalLevel and then calls os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		Init(DefaultConfig())
	}
	globalLogger.Fatal(msg, fields...)
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		Init(DefaultConfig())
	}
	return globalLogger
}
