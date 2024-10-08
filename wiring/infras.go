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
	Logger        logging.Logger
	SnowflakeNode int64
	TokenEngine   token.Engine
}

func InitializeInfras(config config.Config) (Infras, error) {
	infras := Infras{}

	// Logger
	infras.Logger = logging.NewSLogger(logging.Level(config.Server.LogLevel))

	// Snowflake node
	infras.SnowflakeNode = int64(config.Server.NodeID)

	// Token engine
	tokenEngine := token.NewJWTEngine()

	authSecrets := config.Secret.Authentication
	if authSecrets.TokenRSAPrivateKey != "" && authSecrets.TokenRSAPublicKey != "" {
		err := tokenEngine.WithRSA(authSecrets.TokenRSAPrivateKey, authSecrets.TokenRSAPublicKey)
		if err != nil {
			return infras, err
		}
	}

	if authSecrets.TokenHMACSecretKey != "" {
		if err := tokenEngine.WithHMAC(authSecrets.TokenHMACSecretKey); err != nil {
			return infras, err
		}
	}

	infras.TokenEngine = tokenEngine

	return infras, nil
}

func WithInfras(ctx context.Context, infras Infras) context.Context {
	ctx = xcontext.WithLogger(ctx, infras.Logger)
	return ctx
}

func (infras *Infras) NewSnowflakeNode() *snowflake.Node {
	result, err := snowflake.NewNode(infras.SnowflakeNode)
	if err != nil {
		panic(err)
	}
	return result
}
