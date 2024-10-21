package wiring

import (
	"context"

	"github.com/xybor-x/snowflake"
	config "github.com/xybor/todennus-config"
	"github.com/xybor/x/logging"
	"github.com/xybor/x/session"
	"github.com/xybor/x/token"
	"github.com/xybor/x/xcontext"
)

type Infras struct {
	Logger         logging.Logger
	SnowflakeNode  int64
	TokenEngine    token.Engine
	SessionManager *session.Manager
}

func InitializeInfras(config *config.Config) (*Infras, error) {
	infras := &Infras{}

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
	infras.SessionManager = session.NewManager("/", config.Variable.Session.Expiration)

	return infras, nil
}

func WithInfras(ctx context.Context, infras *Infras) context.Context {
	ctx = xcontext.WithLogger(ctx, infras.Logger)
	ctx = xcontext.WithSessionManager(ctx, infras.SessionManager)
	return ctx
}

func (infras *Infras) NewSnowflakeNode() *snowflake.Node {
	result, err := snowflake.NewNode(infras.SnowflakeNode)
	if err != nil {
		panic(err)
	}
	return result
}
