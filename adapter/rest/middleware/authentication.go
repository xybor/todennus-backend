package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/xybor/todennus-backend/config"
	"github.com/xybor/todennus-backend/pkg/token"
	"github.com/xybor/todennus-backend/pkg/xcontext"
	"github.com/xybor/todennus-backend/usecase/dto"
)

func Authentication(engine token.Engine, adminCfg config.AdminSecret) func(http.Handler) http.Handler {
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
						ctx = xcontext.WithRequestUserID(ctx, accessToken.To().Metadata.Subject)
					}
				}
			}

			// TODO: Do not pass the admin secret as token into the
			// authorization header, it is not replay-attack resistence.
			//
			// We need to sign a nonce value by the secret and the backend can
			// check the signature to determine if it is admin.
			//
			// The nonce value can be time unix after that this token will be
			// invalid. The maximum nonce value is now+max_expiration.
			//
			// We can use RSA key to enhance the security of this process.
			adminToken := r.Header.Get("X-ADMIN-AUTHORIZATION")
			if adminToken == adminCfg.SecretKey && adminCfg.SecretKey != "" {
				expiresAt := time.Now().Add(time.Second * time.Duration(adminCfg.MaxExpiration))
				ctx = xcontext.WithAdminExpiresAt(ctx, expiresAt)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
