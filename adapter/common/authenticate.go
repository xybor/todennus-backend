package common

import (
	"context"
	"strings"

	"github.com/xybor/todennus-backend/usecase/dto"
	"github.com/xybor/x/token"
	"github.com/xybor/x/xcontext"
)

func WithAuthenticate(ctx context.Context, authorization string, engine token.Engine) context.Context {
	if authorization == "" {
		return ctx
	}

	tokenType, token, found := strings.Cut(authorization, " ")
	if !found {
		return ctx
	}

	if engine.Type() != tokenType {
		return ctx
	}

	accessToken := dto.OAuth2AccessToken{}
	ok, err := engine.Validate(ctx, token, &accessToken)
	if err != nil {
		xcontext.Logger(ctx).Debug("failed-to-parse-token", "err", err)
		return ctx
	}

	if !ok {
		xcontext.Logger(ctx).Debug("expired token")
		return ctx
	}

	dtoken, err := accessToken.To()
	if err != nil {
		xcontext.Logger(ctx).Warn("failed-to-convert-to-domain-token", "err", err)
		return ctx
	}

	ctx = xcontext.WithRequestUserID(ctx, dtoken.Metadata.Subject)
	ctx = xcontext.WithScope(ctx, dtoken.Scope)

	xcontext.Logger(ctx).Debug("auth-info", "uid", dtoken.Metadata.Subject, "scope", dtoken.Scope.String())

	return ctx
}
