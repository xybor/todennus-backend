package middleware

import (
	"net/http"

	"github.com/xybor/todennus-backend/wiring"
	config "github.com/xybor/todennus-config"
	"github.com/xybor/x/xcontext"
)

func WithInfras(config *config.Config, infras *wiring.Infras) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = wiring.WithInfras(ctx, infras)

			logger := xcontext.Logger(ctx).With("request_id", xcontext.RequestID(ctx)).
				With("node_id", config.Server.NodeID)
			ctx = xcontext.WithLogger(ctx, logger)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
