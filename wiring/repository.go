package wiring

import (
	"context"

	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/infras/database/model"
	"github.com/xybor/todennus-backend/infras/database/redis"
	"github.com/xybor/todennus-backend/usecase/abstraction"
	config "github.com/xybor/todennus-config"
	"github.com/xybor/x/session"
	"github.com/xybor/x/xcrypto"
)

type Repositories struct {
	abstraction.UserRepository
	abstraction.RefreshTokenRepository
	abstraction.OAuth2ClientRepository
	abstraction.SessionRepository
	abstraction.OAuth2AuthorizationCodeRepository
}

func InitializeRepositories(ctx context.Context, config *config.Config, db *Databases) (*Repositories, error) {
	r := &Repositories{}

	r.UserRepository = database.NewUserRepository(db.GormPostgres)
	r.RefreshTokenRepository = database.NewRefreshTokenRepository(db.GormPostgres)
	r.OAuth2ClientRepository = database.NewOAuth2ClientRepository(db.GormPostgres)
	r.SessionRepository = database.NewSessionRepository(
		session.NewCookieStore[model.SessionModel](
			[]byte(config.Secret.Session.AuthenticationKey),
			xcrypto.GenerateAESKeyFromPassword(config.Secret.Session.EncryptionKey, 32),
		))
	r.OAuth2AuthorizationCodeRepository = redis.NewOAuth2AuthorizationCodeRepository(db.Redis)

	return r, nil
}
