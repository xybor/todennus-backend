package middleware

import (
	"net/http"
	"time"

	"github.com/xybor/todennus-backend/pkg/xcontext"
)

func Time(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = xcontext.WithRequestTime(ctx, time.Now())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
