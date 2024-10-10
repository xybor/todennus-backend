package xcontext

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/xybor/todennus-backend/pkg/logging"
	"github.com/xybor/todennus-backend/pkg/scope"
)

type contextKey int

const (
	loggerKey contextKey = iota
	requestTimeKey
	requestUserIDKey
	adminExpiresAtKey
	scopeKey
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

func WithRequestUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, requestUserIDKey, userID)
}

func RequestUserID(ctx context.Context) int64 {
	if val := ctx.Value(requestUserIDKey); val != nil {
		return val.(int64)
	}

	return 0
}

func WithAdminExpiresAt(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, adminExpiresAtKey, t)
}

func AdminExpiresAt(ctx context.Context) time.Time {
	if val := ctx.Value(adminExpiresAtKey); val != nil {
		return val.(time.Time)
	}

	return time.Time{}
}

func IsAdmin(ctx context.Context) bool {
	return AdminExpiresAt(ctx).After(time.Now())
}

func WithScope(ctx context.Context, scopes scope.Scopes) context.Context {
	return context.WithValue(ctx, scopeKey, scopes)
}

func Scope(ctx context.Context) scope.Scopes {
	if val := ctx.Value(scopeKey); val != nil {
		return val.(scope.Scopes)
	}

	return nil
}
