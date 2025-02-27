package logger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name string
		opts LoggerOptions
		want func(*testing.T, *Logger)
	}{
		{
			name: "NoOp logger",
			opts: LoggerOptions{
				NoOp: true,
			},
			want: func(t *testing.T, l *Logger) {
				assert.NotNil(t, l)
				assert.NotNil(t, l.log)
			},
		},
		{
			name: "Debug logger with console encoding",
			opts: LoggerOptions{
				Debug: true,
				JSON:  false,
			},
			want: func(t *testing.T, l *Logger) {
				assert.NotNil(t, l)
				assert.NotNil(t, l.log)
			},
		},
		{
			name: "Production logger with JSON encoding",
			opts: LoggerOptions{
				JSON: true,
			},
			want: func(t *testing.T, l *Logger) {
				assert.NotNil(t, l)
				assert.NotNil(t, l.log)
			},
		},
		{
			name: "Logger with custom time encoder",
			opts: LoggerOptions{
				EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
					enc.AppendString(t.Format(time.RFC3339))
				},
			},
			want: func(t *testing.T, l *Logger) {
				assert.NotNil(t, l)
				assert.NotNil(t, l.log)
			},
		},
		{
			name: "Logger with caller and stack trace",
			opts: LoggerOptions{
				Caller: true,
				Stack:  true,
			},
			want: func(t *testing.T, l *Logger) {
				assert.NotNil(t, l)
				assert.NotNil(t, l.log)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.opts)
			tt.want(t, logger)
		})
	}
}
