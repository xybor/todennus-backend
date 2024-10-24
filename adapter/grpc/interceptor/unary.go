package interceptor

import (
	"context"
	"time"

	"github.com/xybor/todennus-backend/adapter/common"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/todennus-backend/wiring"
	config "github.com/xybor/todennus-config"
	"github.com/xybor/x/token"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xcrypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryInterceptor(config *config.Config, infras *wiring.Infras) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = wiring.WithInfras(ctx, infras)
		ctx = withRequestID(ctx)
		ctx, cancel := withTimeout(ctx, config)
		defer cancel()

		ctx = withAuthenticate(ctx, infras.TokenEngine)

		start := time.Now()
		xcontext.Logger(ctx).Debug("rpc_request", "function", info.FullMethod, "node_id", config.Variable.Server.NodeID)
		resp, err := handler(ctx, req)
		xcontext.Logger(ctx).Debug("rpc_response", "rtt", time.Since(start))

		return resp, err
	}
}

func withRequestID(ctx context.Context) context.Context {
	ctx = xcontext.WithRequestID(ctx, xcrypto.RandString(16))
	logger := xcontext.Logger(ctx).With("request_id", xcontext.RequestID(ctx))
	ctx = xcontext.WithLogger(ctx, logger)
	return ctx
}

func withTimeout(ctx context.Context, config *config.Config) (context.Context, context.CancelFunc) {
	ctx = context.WithoutCancel(ctx)
	timeout := time.Duration(config.Variable.Server.RequestTimeout) * time.Millisecond
	return context.WithTimeoutCause(ctx, timeout, usecase.ErrServerTimeout)
}

func withAuthenticate(ctx context.Context, engine token.Engine) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		xcontext.Logger(ctx).Debug("not-found-metadata")
		return ctx
	}

	authorization := md["authorization"]
	if len(authorization) != 1 {
		xcontext.Logger(ctx).Debug("invalid-or-not-found-authorization-metadata")
		return ctx
	}

	return common.WithAuthenticate(ctx, authorization[0], engine)
}
