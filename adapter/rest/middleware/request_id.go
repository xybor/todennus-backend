package middleware

import (
	"net/http"

	config "github.com/xybor/todennus-config"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xcrypto"
)

func WithRequestID(config config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = xcontext.WithRequestID(ctx, xcrypto.RandString(16))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
