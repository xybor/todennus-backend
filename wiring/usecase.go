package wiring

import (
	"context"
	"time"

	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/usecase"
	config "github.com/xybor/todennus-config"
	"github.com/xybor/x/lock"
)

type Usecases struct {
	abstraction.UserUsecase
	abstraction.OAuth2Usecase
	abstraction.OAuth2ClientUsecase
}

func InitializeUsecases(
	ctx context.Context,
	config *config.Config,
	infras *Infras,
	databases *Databases,
	domains *Domains,
	repositories *Repositories,
) (*Usecases, error) {
	uc := &Usecases{}

	uc.UserUsecase = usecase.NewUserUsecase(
		lock.NewRedisLock(databases.Redis, "user-lock", 10*time.Second),
		repositories.UserRepository,
		domains.UserDomain,
	)

	uc.OAuth2Usecase = usecase.NewOAuth2Usecase(
		infras.TokenEngine,
		config.Variable.OAuth2.IdPLoginURL,
		config.Secret.OAuth2.IdPSecret,
		domains.UserDomain,
		domains.OAuth2FlowDomain,
		domains.OAuth2ClientDomain,
		domains.OAuth2ConsentDomain,
		repositories.UserRepository,
		repositories.RefreshTokenRepository,
		repositories.OAuth2ClientRepository,
		repositories.SessionRepository,
		repositories.OAuth2AuthorizationCodeRepository,
		repositories.OAuth2ConsentRepository,
	)

	uc.OAuth2ClientUsecase = usecase.NewOAuth2ClientUsecase(
		lock.NewRedisLock(databases.Redis, "client-lock", 10*time.Second),
		domains.UserDomain,
		domains.OAuth2ClientDomain,
		repositories.UserRepository,
		repositories.OAuth2ClientRepository,
	)

	return uc, nil
}
