package logger

// LoggerOptions defines configurable parameters for the logger
type LoggerOptions struct {
	Debug  bool // Enable debug logging
	NoOp   bool // Disable logging completely
	JSON   bool // Use JSON format for logs
	Caller bool // Include caller info
	Stack  bool // Include stack trace in debug mode
}

// DefaultOptions returns a default logger configuration
func DefaultOptions() LoggerOptions {
	return LoggerOptions{
		Debug:  false,
		NoOp:   false,
		JSON:   true,
		Caller: true,
		Stack:  false,
	}
}
