package middleware

import (
	"net/http"

	"github.com/xybor/x/session"
	"github.com/xybor/x/xcontext"
)

func WithSession(manager *session.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			session, err := manager.Get(r)
			if err != nil {
				xcontext.Logger(ctx).Debug("failed-to-get-cookie", "err", err)
			} else {
				ctx = xcontext.WithSession(ctx, session)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
