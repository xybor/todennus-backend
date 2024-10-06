package logging

import (
	"github.com/xybor/todennus-backend/pkg/xerror"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelCritical
)

type Logger interface {
	Log(level Level, msg string, a ...any)
	Debug(msg string, a ...any)
	Info(msg string, a ...any)
	Warn(msg string, a ...any)
	Critical(msg string, a ...any)
}

func Serverity2Level(s xerror.Serverity) Level {
	switch s {
	case xerror.ServerityDebug:
		return LevelDebug
	case xerror.ServerityInfo:
		return LevelInfo
	case xerror.ServerityWarn:
		return LevelWarn
	case xerror.ServerityCritical:
		return LevelCritical
	default:
		panic("invalid serverity")
	}
}

func LogServiceError(logger Logger, err xerror.ServiceError, a ...any) {
	logger.Log(Serverity2Level(err.Serverity), err.Error(), a...)
}
