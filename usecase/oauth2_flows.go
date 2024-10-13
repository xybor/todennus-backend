package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/pkg/scope"
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

	userDomain         abstraction.UserDomain
	oauth2ClientDomain abstraction.OAuth2ClientDomain
	oauth2FlowDomain   abstraction.OAuth2FlowDomain

	userRepo         abstraction.UserRepository
	refreshTokenRepo abstraction.RefreshTokenRepository
	oauth2ClientRepo abstraction.OAuth2ClientRepository
}

func NewOAuth2Usecase(
	tokenEngine token.Engine,
	userDomain abstraction.UserDomain,
	oauth2FlowDomain abstraction.OAuth2FlowDomain,
	oauth2ClientDomain abstraction.OAuth2ClientDomain,
	userRepo abstraction.UserRepository,
	refreshTokenRepo abstraction.RefreshTokenRepository,
	oauth2ClientRepo abstraction.OAuth2ClientRepository,
) *OAuth2Usecase {
	return &OAuth2Usecase{
		tokenEngine:        tokenEngine,
		userDomain:         userDomain,
		oauth2FlowDomain:   oauth2FlowDomain,
		oauth2ClientDomain: oauth2ClientDomain,

		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		oauth2ClientRepo: oauth2ClientRepo,
	}
}

func (usecase *OAuth2Usecase) Token(ctx context.Context, req dto.OAuth2TokenRequestDTO) (dto.OAuth2TokenResponseDTO, error) {
	client, err := usecase.oauth2ClientRepo.GetByID(ctx, req.ClientID)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.OAuth2TokenResponseDTO{}, xerror.WrapDebug(ErrClientInvalid).
				WithMessage("client is not found")
		}

		return dto.OAuth2TokenResponseDTO{}, wrapNonDomainError(xerror.ServerityCritical, err)
	}

	switch req.GrantType {
	case GrantTypePassword:
		return usecase.handlePasswordFlow(ctx, req, client)
	case GrantTypeRefreshToken:
		return usecase.handleRefreshTokenFlow(ctx, req, client)
	default:
		return dto.OAuth2TokenResponseDTO{}, fmt.Errorf("%w: %s", ErrGrantTypeInvalid, req.GrantType)
	}
}

func (usecase *OAuth2Usecase) handlePasswordFlow(
	ctx context.Context,
	req dto.OAuth2TokenRequestDTO,
	client domain.OAuth2Client,
) (dto.OAuth2TokenResponseDTO, error) {
	err := usecase.oauth2ClientDomain.ValidateClient(
		client, req.ClientID, req.ClientSecret, domain.RequireConfidential)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, wrapDomainError(err)
	}

	// Get the user information.
	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.OAuth2TokenResponseDTO{}, xerror.WrapDebug(ErrUsernamePasswordInvalid)
		}

		return dto.OAuth2TokenResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	// Validate password.
	ok, err := usecase.userDomain.Validate(user.HashedPass, req.Password)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, wrapDomainError(err)
	}
	if !ok {
		return dto.OAuth2TokenResponseDTO{}, xerror.WrapDebug(ErrUsernamePasswordInvalid)
	}

	// Validate scope.
	requestedScope := domain.ScopeEngine.ParseScopes(req.Scope)
	if !requestedScope.LessThan(client.AllowedScope) {
		return dto.OAuth2TokenResponseDTO{}, xerror.WrapDebug(ErrScopeInvalid).
			WithMessage("client has not permission to request this scope")
	}

	finalScope := requestedScope.Intersect(user.AllowedScope)

	// Generate both tokens.
	accessToken, refreshToken, err := usecase.generateAccessAndRefreshTokens(finalScope, user)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, err
	}

	// Serialize both tokens.
	accessTokenString, refreshTokenString, err := usecase.serializeAccessAndRefreshTokens(ctx, accessToken, refreshToken)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, err
	}

	// Store refresh token information.
	err = usecase.refreshTokenRepo.Create(ctx, refreshToken.Metadata.Id, accessToken.Metadata.Id, 0)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.OAuth2TokenResponseDTO{
		AccessToken:  accessTokenString,
		TokenType:    usecase.tokenEngine.Type(),
		ExpiresIn:    accessToken.Metadata.ExpiresIn,
		RefreshToken: refreshTokenString,
		Scope:        finalScope.String(),
	}, nil
}

func (usecase *OAuth2Usecase) handleRefreshTokenFlow(
	ctx context.Context,
	req dto.OAuth2TokenRequestDTO,
	client domain.OAuth2Client,
) (dto.OAuth2TokenResponseDTO, error) {
	err := usecase.oauth2ClientDomain.ValidateClient(
		client, req.ClientID, req.ClientSecret, domain.DependOnClientConfidential)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, wrapDomainError(err)
	}

	// Check the current refresh token
	curRefreshToken := dto.OAuth2RefreshToken{}
	ok, err := usecase.tokenEngine.Validate(ctx, req.RefreshToken, &curRefreshToken)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	if !ok {
		return dto.OAuth2TokenResponseDTO{}, xerror.WrapDebug(ErrRefreshTokenInvalid)
	}

	// Generate the next refresh token.
	domainCurRefreshToken, err := curRefreshToken.To()
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	refreshToken, err := usecase.oauth2FlowDomain.NextRefreshToken(domainCurRefreshToken)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, wrapDomainError(err)
	}

	// Get the user.
	user, err := usecase.userRepo.GetByID(ctx, refreshToken.Metadata.Subject)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	// Generate access token.
	accessToken, err := usecase.oauth2FlowDomain.CreateAccessToken(
		domainCurRefreshToken.Metadata.Audience, domainCurRefreshToken.Scope, user)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, wrapDomainError(err)
	}

	// Serialize both tokens.
	accessTokenString, refreshTokenString, err := usecase.serializeAccessAndRefreshTokens(ctx, accessToken, refreshToken)
	if err != nil {
		return dto.OAuth2TokenResponseDTO{}, err
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

			return dto.OAuth2TokenResponseDTO{}, xerror.WrapDebug(ErrRefreshTokenStolen)
		}

		return dto.OAuth2TokenResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.OAuth2TokenResponseDTO{
		AccessToken:  accessTokenString,
		TokenType:    usecase.tokenEngine.Type(),
		ExpiresIn:    accessToken.Metadata.ExpiresIn,
		RefreshToken: refreshTokenString,
		Scope:        domainCurRefreshToken.Scope.String(),
	}, nil
}

func (usecase *OAuth2Usecase) generateAccessAndRefreshTokens(
	scope scope.Scopes,
	user domain.User,
) (domain.OAuth2AccessToken, domain.OAuth2RefreshToken, error) {
	accessToken, err := usecase.oauth2FlowDomain.CreateAccessToken("", scope, user)
	if err != nil {
		return domain.OAuth2AccessToken{}, domain.OAuth2RefreshToken{}, wrapDomainError(err)
	}

	refreshToken, err := usecase.oauth2FlowDomain.CreateRefreshToken("", scope, user.ID)
	if err != nil {
		return domain.OAuth2AccessToken{}, domain.OAuth2RefreshToken{}, wrapDomainError(err)
	}

	return accessToken, refreshToken, nil
}

func (usecase *OAuth2Usecase) serializeAccessAndRefreshTokens(
	ctx context.Context,
	accessToken domain.OAuth2AccessToken,
	refreshToken domain.OAuth2RefreshToken,
) (string, string, error) {
	accessTokenString, err := usecase.tokenEngine.Generate(ctx, dto.OAuth2AccessTokenFromDomain(accessToken))
	if err != nil {
		return "", "", wrapNonDomainError(xerror.ServerityDebug, err)
	}

	refreshTokenString, err := usecase.tokenEngine.Generate(ctx, dto.OAuth2RefreshTokenFromDomain(refreshToken))
	if err != nil {
		return "", "", wrapNonDomainError(xerror.ServerityDebug, err)
	}

	return accessTokenString, refreshTokenString, nil
}
