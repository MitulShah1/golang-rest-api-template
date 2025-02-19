package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the logger with given options
func Init(opts LoggerOptions) {
	if opts.NoOp {
		log = zap.NewNop() // NoOp logger when logging is disabled
		return
	}

	cfg := zap.NewProductionConfig()

	// Set log format
	if !opts.JSON {
		cfg.Encoding = "console"
	}

	// Configure logging level
	if opts.Debug {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Enable caller and stack trace if configured
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-readable timestamp
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.StacktraceKey = ""

	if opts.Caller {
		cfg.EncoderConfig.CallerKey = "caller"
	}

	if opts.Stack {
		cfg.EncoderConfig.StacktraceKey = "stacktrace"
	}

	// Build logger
	var err error
	log, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// Info logs an informational message with key-value pairs
func Info(msg string, keysAndValues ...interface{}) {
	log.Sugar().Infow(msg, keysAndValues...)
}

// Debug logs a debug message with key-value pairs
func Debug(msg string, keysAndValues ...interface{}) {
	log.Sugar().Debugw(msg, keysAndValues...)
}

// Error logs an error message with key-value pairs
func Error(msg string, keysAndValues ...interface{}) {
	log.Sugar().Errorw(msg, keysAndValues...)
}

// Warn logs a warning message with key-value pairs
func Warn(msg string, keysAndValues ...interface{}) {
	log.Sugar().Warnw(msg, keysAndValues...)
}

// Sync flushes any buffered log entries
func Sync() {
	_ = log.Sync()
}
