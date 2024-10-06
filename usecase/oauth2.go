package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/pkg/token"
	"github.com/xybor/todennus-backend/pkg/xerror"
	"github.com/xybor/todennus-backend/usecase/abstraction"
	"github.com/xybor/todennus-backend/usecase/dto"
)

const (
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypePassword          = "password"
	GrantTypeClientCredentials = "client_credentials"
	GrantTypeRefreshToken      = "refresh_token"

	// TODO: Support later
	GrantTypeDevice = "urn:ietf:params:oauth:grant-type:device_code"
)

type OAuth2Usecase struct {
	tokenEngine token.Engine

	userDomain   abstraction.UserDomain
	oauth2Domain abstraction.OAuth2Domain

	userRepo abstraction.UserRepository
}

func NewOAuth2Usecase(
	tokenEngine token.Engine,
	userDomain abstraction.UserDomain,
	oauth2Domain abstraction.OAuth2Domain,
	userRepo abstraction.UserRepository,
) *OAuth2Usecase {
	return &OAuth2Usecase{
		tokenEngine:  tokenEngine,
		userDomain:   userDomain,
		oauth2Domain: oauth2Domain,
		userRepo:     userRepo,
	}
}

func (usecase *OAuth2Usecase) Token(ctx context.Context, req dto.OAuth2TokenRequest) (dto.OAuth2TokenResponse, error) {
	switch req.GrantType {
	case GrantTypePassword:
		return usecase.handlePasswordFlow(ctx, req)
	default:
		return dto.OAuth2TokenResponse{}, fmt.Errorf("%w: %s", ErrInvalidGrantType, req.GrantType)
	}
}

func (usecase *OAuth2Usecase) handlePasswordFlow(
	ctx context.Context,
	req dto.OAuth2TokenRequest,
) (dto.OAuth2TokenResponse, error) {
	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.OAuth2TokenResponse{}, xerror.WrapDebug(ErrUsernamePasswordInvalid)
		}

		return dto.OAuth2TokenResponse{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	ok, err := usecase.userDomain.Validate(user.HashedPass, req.Password)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapDomainError(err)
	}

	if !ok {
		return dto.OAuth2TokenResponse{}, xerror.WrapDebug(ErrUsernamePasswordInvalid)
	}

	accessToken, err := usecase.oauth2Domain.CreateAccessToken("", user.ID)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapDomainError(err)
	}

	accessTokenString, err := usecase.tokenEngine.Generate(
		ctx, dto.OAuth2AccessTokenFromDomain(accessToken))
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapNonDomainError(xerror.ServerityCritical, err)
	}

	return dto.OAuth2TokenResponse{
		AccessToken: accessTokenString,
		TokenType:   usecase.tokenEngine.Type(),
		ExpiresIn:   accessToken.Metadata.ExpiresIn,
	}, nil
}
