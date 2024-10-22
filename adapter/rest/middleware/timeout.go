package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/xybor/todennus-backend/usecase"
	config "github.com/xybor/todennus-config"
)

func Timeout(config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithoutCancel(ctx)

			timeout := time.Duration(config.Variable.Server.RequestTimeout) * time.Millisecond
			ctx, cancel := context.WithTimeoutCause(ctx, timeout, usecase.ErrServerTimeout)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
