package token

import "context"

type Claims interface {
	Valid() error
}

type Engine interface {
	Type() string
	Generate(ctx context.Context, inner Claims) (string, error)
	Validate(ctx context.Context, token string, claims Claims) (bool, error)
}
