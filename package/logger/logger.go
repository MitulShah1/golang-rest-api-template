package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

// Init initializes the logger with given options
func NewLogger(opts LoggerOptions) *Logger {
	logger := &Logger{}

	if opts.NoOp {
		logger.log = zap.NewNop() // NoOp logger when logging is disabled
		return logger
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
	if opts.EncodeTime != nil {
		cfg.EncoderConfig.EncodeTime = opts.EncodeTime
	} else {
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-readable timestamp
	}
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
	logger.log, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	return logger
}

// Info logs an informational message with key-value pairs
func (lg Logger) Info(msg string, keysAndValues ...interface{}) {
	lg.log.Sugar().Infow(msg, keysAndValues...)
}

// Debug logs a debug message with key-value pairs
func (lg Logger) Debug(msg string, keysAndValues ...interface{}) {
	lg.log.Sugar().Debugw(msg, keysAndValues...)
}

// Error logs an error message with key-value pairs
func (lg Logger) Error(msg string, keysAndValues ...interface{}) {
	lg.log.Sugar().Errorw(msg, keysAndValues...)
}

// Warn logs a warning message with key-value pairs
func (lg Logger) Warn(msg string, keysAndValues ...interface{}) {
	lg.log.Sugar().Warnw(msg, keysAndValues...)
}

// Fatal logs a fatal message with key-value pairs
func (lg Logger) Fatal(msg string, keysAndValues ...interface{}) {
	lg.log.Sugar().Fatalw(msg, keysAndValues...)
}

// Sync flushes any buffered log entries
func (lg Logger) Sync() {
	_ = lg.log.Sync()
}
