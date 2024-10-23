package wiring

import (
	"context"
	"time"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/usecase/abstraction"
	config "github.com/xybor/todennus-config"
)

type Domains struct {
	abstraction.UserDomain
	abstraction.OAuth2FlowDomain
	abstraction.OAuth2ClientDomain
	abstraction.OAuth2ConsentDomain
}

func InitializeDomains(ctx context.Context, config *config.Config, infras *Infras) (*Domains, error) {
	var err error
	domains := &Domains{}

	domains.UserDomain, err = domain.NewUserDomain(infras.NewSnowflakeNode())
	if err != nil {
		return nil, err
	}

	domains.OAuth2FlowDomain, err = domain.NewOAuth2FlowDomain(
		infras.NewSnowflakeNode(),
		config.Variable.Authentication.TokenIssuer,
		time.Duration(config.Variable.OAuth2.AuthorizationCodeFlowExpiration)*time.Second,
		time.Duration(config.Variable.OAuth2.AuthenticationCallbackExpiration)*time.Second,
		time.Duration(config.Variable.OAuth2.SessionUpdateExpiration)*time.Second,
		time.Duration(config.Variable.Session.Expiration)*time.Second,
		time.Duration(config.Variable.Authentication.AccessTokenExpiration)*time.Second,
		time.Duration(config.Variable.Authentication.RefreshTokenExpiration)*time.Second,
		time.Duration(config.Variable.Authentication.IDTokenExpiration)*time.Second,
	)
	if err != nil {
		return nil, err
	}

	domains.OAuth2ClientDomain, err = domain.NewOAuth2ClientDomain(
		infras.NewSnowflakeNode(),
		config.Variable.OAuth2.ClientSecretLength,
	)
	if err != nil {
		return nil, err
	}

	domains.OAuth2ConsentDomain = domain.NewOAuth2ConsentDomain(
		time.Duration(config.Variable.OAuth2.ConsentSessionExpiration)*time.Second,
		time.Duration(config.Variable.OAuth2.ConsentExpiration)*time.Second,
	)

	return domains, nil
}
