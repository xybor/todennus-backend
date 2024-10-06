package middleware

import (
	"net/http"

	"github.com/xybor/todennus-backend/pkg/xcontext"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		xcontext.Logger(ctx).Debug(
			r.RequestURI,
			"rip", r.RemoteAddr,
			"request_id", xcontext.RequestID(ctx),
		)

		next.ServeHTTP(w, r)
	})
}
