package interceptor

import (
	"context"
	"time"

	"github.com/xybor/todennus-backend/usecase"
	ucdto "github.com/xybor/todennus-backend/usecase/dto"
	"github.com/xybor/todennus-backend/wiring"
	config "github.com/xybor/todennus-config"
	"github.com/xybor/x/token"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xcrypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnarySetupContext(config *config.Config, infras *wiring.Infras) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = wiring.WithInfras(ctx, infras)
		ctx = xcontext.WithRequestID(ctx, xcrypto.RandString(16))

		logger := xcontext.Logger(ctx).With("request_id", xcontext.RequestID(ctx))
		ctx = xcontext.WithLogger(ctx, logger)
		ctx = context.WithoutCancel(ctx)

		timeout := time.Duration(config.Variable.Server.RequestTimeout) * time.Millisecond
		ctx, cancel := context.WithTimeoutCause(ctx, timeout, usecase.ErrServerTimeout)
		defer cancel()

		start := time.Now()
		xcontext.Logger(ctx).Debug("rpc_request", "function", info.FullMethod, "node_id", config.Variable.Server.NodeID)
		resp, err := handler(ctx, req)
		xcontext.Logger(ctx).Debug("rpc_response", "rtt", time.Since(start))

		return resp, err
	}
}

func UnaryAuthenticate(engine token.Engine) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			authorization := md["authorization"]
			if len(authorization) == 2 {
				tokenType, token := authorization[0], authorization[1]
				if engine.Type() == tokenType {
					accessToken := ucdto.OAuth2AccessToken{}

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

							xcontext.Logger(ctx).Debug("auth-info",
								"user-id", domainAccessToken.Metadata.Subject,
								"scope", domainAccessToken.Scope.String(),
							)
						}
					}
				}
			}
		}

		return handler(ctx, req)
	}
}
