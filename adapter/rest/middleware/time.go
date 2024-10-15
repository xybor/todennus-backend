package middleware

import (
	"net/http"
	"time"

	"github.com/xybor/x/xcontext"
)

func RoundTripTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		xcontext.Logger(r.Context()).Debug("response", "rtt", time.Since(start))
	})
}
