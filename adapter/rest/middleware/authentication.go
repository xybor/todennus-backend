package middleware

import (
	"net/http"
	"strings"

	"github.com/xybor/todennus-backend/adapter/rest/response"
	"github.com/xybor/todennus-backend/adapter/rest/standard"
	"github.com/xybor/todennus-backend/usecase/dto"
	"github.com/xybor/x/token"
	"github.com/xybor/x/xcontext"
)

func Authentication(engine token.Engine) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authenticationHeader := r.Header.Get("Authorization")
			tokenType, token, found := strings.Cut(authenticationHeader, " ")
			if found {
				if engine.Type() == tokenType {
					accessToken := dto.OAuth2AccessToken{}

					ok, err := engine.Validate(ctx, token, &accessToken)
					if err != nil {
						xcontext.Logger(ctx).Debug("failed-to-parse-token", "err", err)
					} else if ok {
						domainAccessToken, err := accessToken.To()
						if err != nil {
							xcontext.Logger(ctx).Warn("failed-to-convert-to-domain-token", "err", err)
						} else {
							ctx = xcontext.WithRequestUserID(ctx, domainAccessToken.Metadata.Subject)
							ctx = xcontext.WithScope(ctx, domainAccessToken.Scope)

							xcontext.Logger(ctx).Debug("request-auth",
								"user-id", domainAccessToken.Metadata.Subject,
								"scope", domainAccessToken.Scope.String(),
							)
						}
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuthentication(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if xcontext.RequestUserID(ctx) == 0 {
			response.Write(ctx, w, http.StatusUnauthorized,
				standard.NewErrorResponseWithMessage(ctx, "unauthenticated", "require authentication to access api"))
		} else {
			handler(w, r)
		}
	}
}
