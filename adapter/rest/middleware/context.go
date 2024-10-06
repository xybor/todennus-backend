package middleware

import (
	"net/http"

	"github.com/xybor/todennus-backend/wiring"
)

func WithInfras(infras wiring.Infras) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = wiring.WithInfras(ctx, infras)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
