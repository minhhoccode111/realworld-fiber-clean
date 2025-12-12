package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Interface -.
// Interface is an interface for logging.
type Interface interface {
	Debug(message any, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message any, args ...any)
	Fatal(message any, args ...any)
}

// Logger -.
// Logger implements the Interface for logging.
type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
// New creates a new logger with the specified level.
func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{
		logger: &logger,
	}
}

// Debug -.
// Debug logs a debug message.
func (l *Logger) Debug(message any, args ...any) {
	l.msg("debug", message, args...)
}

// Info -.
// Info logs an info message.
func (l *Logger) Info(message string, args ...any) {
	l.log(message, args...)
}

// Warn -.
// Warn logs a warning message.
func (l *Logger) Warn(message string, args ...any) {
	l.log(message, args...)
}

// Error -.
// Error logs an error message.
func (l *Logger) Error(message any, args ...any) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, args...)
	}

	l.msg("error", message, args...)
}

// Fatal -.
// Fatal logs a fatal message and exits.
func (l *Logger) Fatal(message any, args ...any) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}

func (l *Logger) log(message string, args ...any) {
	if len(args) == 0 {
		l.logger.Info().Msg(message)
	} else {
		l.logger.Info().Msgf(message, args...)
	}
}

func (l *Logger) msg(level string, message any, args ...any) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
