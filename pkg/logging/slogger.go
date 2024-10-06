package logging

import (
	"context"
	"log/slog"
	"os"
)

var _ Logger = (*SLogger)(nil)

type SLogger struct {
	core *slog.Logger
}

func LevelToSlogLevel(level Level) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelCritical:
		return slog.LevelError
	default:
		panic("invalid level")
	}
}

func NewSLogger(level Level) *SLogger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: LevelToSlogLevel(level),
	})
	core := slog.New(handler)
	if !core.Enabled(context.Background(), LevelToSlogLevel(level)) {
		panic("Cannot enable level")
	}
	return &SLogger{core: core}
}

func (l *SLogger) Log(level Level, msg string, a ...any) {
	l.core.Log(context.Background(), LevelToSlogLevel(level), msg, a...)
}

func (l *SLogger) Debug(msg string, a ...any) {
	l.Log(LevelDebug, msg, a...)
}

func (l *SLogger) Info(msg string, a ...any) {
	l.Log(LevelInfo, msg, a...)
}

func (l *SLogger) Warn(msg string, a ...any) {
	l.Log(LevelWarn, msg, a...)
}

func (l *SLogger) Critical(msg string, a ...any) {
	l.Log(LevelCritical, msg, a...)
}
