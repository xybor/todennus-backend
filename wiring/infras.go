package wiring

import (
	"context"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/config"
	"github.com/xybor/todennus-backend/pkg/logging"
	"github.com/xybor/todennus-backend/pkg/token"
	"github.com/xybor/todennus-backend/pkg/xcontext"
)

type Infras struct {
	Logger      logging.Logger
	Snowflake   *snowflake.Node
	TokenEngine token.Engine
}

func InitializeInfras(config config.Config) (Infras, error) {
	ctxConfig := Infras{}

	var err error

	// Logger
	ctxConfig.Logger = logging.NewSLogger(logging.Level(config.Server.LogLevel))

	// Snowflake node
	ctxConfig.Snowflake, err = snowflake.NewNode(int64(config.Server.NodeID))
	if err != nil {
		return ctxConfig, err
	}

	// Token engine
	tokenEngine := token.NewJWTEngine()

	authSecrets := config.Secret.Authentication
	if authSecrets.TokenRSAPrivateKey != "" && authSecrets.TokenRSAPublicKey != "" {
		err := tokenEngine.WithRSA(authSecrets.TokenRSAPrivateKey, authSecrets.TokenRSAPublicKey)
		if err != nil {
			return ctxConfig, err
		}
	}

	if authSecrets.TokenHMACSecretKey != "" {
		if err := tokenEngine.WithHMAC(authSecrets.TokenHMACSecretKey); err != nil {
			return ctxConfig, err
		}
	}

	ctxConfig.TokenEngine = tokenEngine

	return ctxConfig, nil
}

func WithInfras(ctx context.Context, infras Infras) context.Context {
	ctx = xcontext.WithLogger(ctx, infras.Logger)
	return ctx
}
