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
	abstraction.OAuth2FlowDomain
	abstraction.OAuth2ClientDomain
}

func InitializeDomains(ctx context.Context, config config.Config, infras Infras) (Domains, error) {
	domains := Domains{}

	var finalErr error
	var err error

	domains.UserDomain, err = domain.NewUserDomain(infras.NewSnowflakeNode())
	finalErr = errors.Join(finalErr, err)

	domains.OAuth2FlowDomain, err = domain.NewOAuth2FlowDomain(
		infras.NewSnowflakeNode(),
		config.Variable.Authentication.TokenIssuer,
		time.Duration(config.Variable.Authentication.AccessTokenExpiration)*time.Second,
		time.Duration(config.Variable.Authentication.RefreshTokenExpiration)*time.Second,
		time.Duration(config.Variable.Authentication.IDTokenExpiration)*time.Second,
	)
	finalErr = errors.Join(finalErr, err)

	domains.OAuth2ClientDomain, err = domain.NewOAuth2ClientDomain(
		infras.NewSnowflakeNode(),
		config.Variable.OAuth2.ClientSecretLength,
	)
	finalErr = errors.Join(finalErr, err)

	return domains, finalErr
}
