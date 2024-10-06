package xcontext

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/xybor/todennus-backend/pkg/logging"
)

type contextKey int

const (
	loggerKey contextKey = iota
	requestTimeKey
)

func WithLogger(ctx context.Context, logger logging.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func Logger(ctx context.Context) logging.Logger {
	return ctx.Value(loggerKey).(logging.Logger)
}

func RequestID(ctx context.Context) string {
	return ctx.Value(middleware.RequestIDKey).(string)
}

func WithRequestTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, requestTimeKey, t)
}

func RequestTime(ctx context.Context) time.Time {
	return ctx.Value(requestTimeKey).(time.Time)
}
