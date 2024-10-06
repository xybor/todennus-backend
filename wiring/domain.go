package wiring

import (
	"context"
	"errors"
	"time"

	"github.com/xybor/todennus-backend/config"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/usecase/abstraction"
)

type Domains struct {
	abstraction.UserDomain
	abstraction.OAuth2Domain
}

func InitializeDomains(ctx context.Context, config config.Config, infras Infras) (Domains, error) {
	domains := Domains{}

	var finalErr error
	var err error

	domains.UserDomain, err = domain.NewUserDomain(infras.NewSnowflakeNode())
	finalErr = errors.Join(finalErr, err)

	domains.OAuth2Domain, err = domain.NewOAuth2Domain(
		infras.NewSnowflakeNode(),
		config.Variable.Authentication.TokenIssuer,
		time.Duration(config.Variable.Authentication.AccessTokenExpiration)*time.Second,
		time.Duration(config.Variable.Authentication.RefreshTokenExpiration)*time.Second,
		time.Duration(config.Variable.Authentication.IDTokenExpiration)*time.Second,
	)
	finalErr = errors.Join(finalErr, err)

	return domains, finalErr
}
