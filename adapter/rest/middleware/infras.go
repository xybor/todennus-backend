package middleware

import (
	"net/http"

	"github.com/xybor/todennus-backend/wiring"
	"github.com/xybor/x/xcontext"
)

func WithInfras(infras wiring.Infras) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = wiring.WithInfras(ctx, infras)
			ctx = xcontext.WithLogger(ctx, xcontext.Logger(ctx).With("request_id", xcontext.RequestID(ctx)))

			xcontext.Logger(ctx).Debug("request", "uri", r.RequestURI, "rip", r.RemoteAddr)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
