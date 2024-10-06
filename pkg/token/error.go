package token

import "errors"

var (
	ErrSigningKeyInvalid            = errors.New("invalid signing key")
	ErrTokenSigningMethodNotSupport = errors.New("not supported signing method")
	ErrTokenExpired                 = errors.New("expired token")
	ErrTokenNotYetValid             = errors.New("not yet valid token")
	ErrTokenInvalidIssuer           = errors.New("token issuer is invalid")
	ErrTokenInvalidFormat           = errors.New("token format is invalid")
)
