package wiring

import (
	"context"

	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/usecase/abstraction"
)

type Repositories struct {
	abstraction.UserRepository
	abstraction.RefreshTokenRepository
	abstraction.OAuth2ClientRepository
}

func InitializeRepositories(ctx context.Context, db Databases) (Repositories, error) {
	r := Repositories{}

	r.UserRepository = database.NewUserRepository(db.GormPostgres)
	r.RefreshTokenRepository = database.NewRefreshTokenRepository(db.GormPostgres)
	r.OAuth2ClientRepository = database.NewOAuth2ClientRepository(db.GormPostgres)

	return r, nil
}
