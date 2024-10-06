package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/pkg/token"
	"github.com/xybor/todennus-backend/pkg/xcontext"
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

	userRepo         abstraction.UserRepository
	refreshTokenRepo abstraction.RefreshTokenRepository
}

func NewOAuth2Usecase(
	tokenEngine token.Engine,
	userDomain abstraction.UserDomain,
	oauth2Domain abstraction.OAuth2Domain,
	userRepo abstraction.UserRepository,
	refreshTokenRepo abstraction.RefreshTokenRepository,
) *OAuth2Usecase {
	return &OAuth2Usecase{
		tokenEngine:      tokenEngine,
		userDomain:       userDomain,
		oauth2Domain:     oauth2Domain,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (usecase *OAuth2Usecase) Token(ctx context.Context, req dto.OAuth2TokenRequest) (dto.OAuth2TokenResponse, error) {
	switch req.GrantType {
	case GrantTypePassword:
		return usecase.handlePasswordFlow(ctx, req)
	case GrantTypeRefreshToken:
		return usecase.handleRefreshTokenFlow(ctx, req)
	default:
		return dto.OAuth2TokenResponse{}, fmt.Errorf("%w: %s", ErrInvalidGrantType, req.GrantType)
	}
}

func (usecase *OAuth2Usecase) handlePasswordFlow(
	ctx context.Context,
	req dto.OAuth2TokenRequest,
) (dto.OAuth2TokenResponse, error) {
	// Get the user information.
	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.OAuth2TokenResponse{}, xerror.WrapDebug(ErrUsernamePasswordInvalid)
		}

		return dto.OAuth2TokenResponse{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	// Validate password.
	ok, err := usecase.userDomain.Validate(user.HashedPass, req.Password)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapDomainError(err)
	}
	if !ok {
		return dto.OAuth2TokenResponse{}, xerror.WrapDebug(ErrUsernamePasswordInvalid)
	}

	// Generate access token.
	accessToken, err := usecase.oauth2Domain.CreateAccessToken("", user.ID)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapDomainError(err)
	}

	refreshToken, err := usecase.oauth2Domain.CreateRefreshToken("", user.ID)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapDomainError(err)
	}

	// Serialize both tokens.
	accessTokenString, refreshTokenString, err := usecase.serializeTokens(ctx, accessToken, refreshToken)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapNonDomainError(xerror.ServerityDebug, err)
	}

	// Store refresh token information.
	err = usecase.refreshTokenRepo.Save(ctx, refreshToken.Metadata.Id, accessToken.Metadata.Id, 0)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.OAuth2TokenResponse{
		AccessToken:  accessTokenString,
		TokenType:    usecase.tokenEngine.Type(),
		ExpiresIn:    accessToken.Metadata.ExpiresIn,
		RefreshToken: refreshTokenString,
	}, nil
}

func (usecase *OAuth2Usecase) handleRefreshTokenFlow(
	ctx context.Context,
	req dto.OAuth2TokenRequest,
) (dto.OAuth2TokenResponse, error) {
	// Check the current refresh token
	curRefreshToken := dto.OAuth2RefreshToken{}
	ok, err := usecase.tokenEngine.Validate(ctx, req.RefreshToken, &curRefreshToken)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	if !ok {
		return dto.OAuth2TokenResponse{}, xerror.WrapDebug(ErrInvalidRefreshToken)
	}

	// Generate the next refresh token.
	domainCurRefreshToken := curRefreshToken.To()
	refreshToken, err := usecase.oauth2Domain.NextRefreshToken(domainCurRefreshToken)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapDomainError(err)
	}

	// Generate access token.
	accessToken, err := usecase.oauth2Domain.CreateAccessToken(
		domainCurRefreshToken.Metadata.Audience,
		domainCurRefreshToken.Metadata.Subject,
	)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapDomainError(err)
	}

	// Serialize both tokens.
	accessTokenString, refreshTokenString, err := usecase.serializeTokens(ctx, accessToken, refreshToken)
	if err != nil {
		return dto.OAuth2TokenResponse{}, wrapNonDomainError(xerror.ServerityDebug, err)
	}

	// Store the seq number again.
	err = usecase.refreshTokenRepo.UpdateByRefreshTokenID(
		ctx,
		domainCurRefreshToken.Metadata.Id,
		accessToken.Metadata.Id,
		domainCurRefreshToken.SequenceNumber,
	)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			err = usecase.refreshTokenRepo.DeleteByRefreshTokenID(ctx, domainCurRefreshToken.Metadata.Id)
			if err != nil {
				xcontext.Logger(ctx).Warn("failed to delete stolen token", "err", err)
			}

			return dto.OAuth2TokenResponse{}, xerror.WrapDebug(ErrStolenRefreshToken)
		}

		return dto.OAuth2TokenResponse{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.OAuth2TokenResponse{
		AccessToken:  accessTokenString,
		TokenType:    usecase.tokenEngine.Type(),
		ExpiresIn:    accessToken.Metadata.ExpiresIn,
		RefreshToken: refreshTokenString,
	}, nil
}

func (usecase *OAuth2Usecase) serializeTokens(
	ctx context.Context,
	accessToken domain.OAuth2AccessToken,
	refreshToken domain.OAuth2RefreshToken,
) (string, string, error) {
	accessTokenString, err := usecase.tokenEngine.Generate(ctx, dto.OAuth2AccessTokenFromDomain(accessToken))
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := usecase.tokenEngine.Generate(ctx, dto.OAuth2RefreshTokenFromDomain(refreshToken))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
