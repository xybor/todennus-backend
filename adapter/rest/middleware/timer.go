package middleware

import (
	"net/http"
	"time"

	"github.com/xybor/x/xcontext"
)

func Timer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		xcontext.Logger(ctx).Debug("request", "uri", r.RequestURI, "rip", r.RemoteAddr)

		start := time.Now()
		next.ServeHTTP(w, r)

		xcontext.Logger(ctx).Debug("response", "rtt", time.Since(start))
	})
}
