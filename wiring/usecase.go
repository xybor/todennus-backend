package wiring

import (
	"context"

	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/usecase"
)

type Usecases struct {
	abstraction.UserUsecase
	abstraction.OAuth2Usecase
}

func InitializeUsecases(
	ctx context.Context,
	infras Infras,
	domains Domains,
	repositories Repositories,
) (Usecases, error) {
	uc := Usecases{}
	uc.UserUsecase = usecase.NewUserUsecase(repositories.UserRepository, domains.UserDomain)
	uc.OAuth2Usecase = usecase.NewOAuth2Usecase(
		infras.TokenEngine, domains.UserDomain, domains.OAuth2Domain, repositories.UserRepository)

	return uc, nil
}
