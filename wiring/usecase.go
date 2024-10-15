package wiring

import (
	"context"
	"time"

	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/x/lock"
)

type Usecases struct {
	abstraction.UserUsecase
	abstraction.OAuth2Usecase
	abstraction.OAuth2ClientUsecase
}

func InitializeUsecases(
	ctx context.Context,
	infras Infras,
	domains Domains,
	repositories Repositories,
) (Usecases, error) {
	uc := Usecases{}

	uc.UserUsecase = usecase.NewUserUsecase(
		lock.NewRedisLock(infras.RedisClient, "user-lock", 10*time.Second),
		repositories.UserRepository,
		domains.UserDomain,
	)

	uc.OAuth2Usecase = usecase.NewOAuth2Usecase(
		infras.TokenEngine,
		domains.UserDomain,
		domains.OAuth2FlowDomain,
		domains.OAuth2ClientDomain,
		repositories.UserRepository,
		repositories.RefreshTokenRepository,
		repositories.OAuth2ClientRepository,
	)

	uc.OAuth2ClientUsecase = usecase.NewOAuth2ClientUsecase(
		lock.NewRedisLock(infras.RedisClient, "client-lock", 10*time.Second),
		domains.UserDomain,
		domains.OAuth2ClientDomain,
		repositories.UserRepository,
		repositories.OAuth2ClientRepository,
	)

	return uc, nil
}
