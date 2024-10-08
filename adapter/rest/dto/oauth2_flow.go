package dto

import (
	"net/http"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/pkg/xerror"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/todennus-backend/usecase/dto"
)

type OAuth2TokenRequestDTO struct {
	GrantType string `form:"grant_type"`

	ClientID     int64  `form:"client_id"`
	ClientSecret string `form:"client_secret"`

	// Authorization Code Flow
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	CodeVerifier string `form:"code_verifier"` // with PKCE

	// Resource Owner Password Credentials Flow
	Username string `form:"username"`
	Password string `form:"password"`

	// Refresh Token Flow
	RefreshToken string `form:"refresh_token"`
}

func (req OAuth2TokenRequestDTO) To() dto.OAuth2TokenRequestDTO {
	return dto.OAuth2TokenRequestDTO{
		GrantType: req.GrantType,

		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,

		Code:         req.Code,
		RedirectURI:  req.RedirectURI,
		CodeVerifier: req.CodeVerifier,

		Username: req.Username,
		Password: req.Password,

		RefreshToken: req.RefreshToken,
	}
}

type OAuth2TokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type OAuth2TokenErrorResponseDTO struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri,omitempty"`
}

func NewOAuth2TokenResponseDTO(resp dto.OAuth2TokenResponseDTO) OAuth2TokenResponseDTO {
	return OAuth2TokenResponseDTO{
		AccessToken:  resp.AccessToken,
		TokenType:    resp.TokenType,
		ExpiresIn:    resp.ExpiresIn,
		RefreshToken: resp.RefreshToken,
		Scope:        resp.Scope,
	}
}

func NewOAuth2TokenErrorResponseDTO(err error) (int, OAuth2TokenErrorResponseDTO) {
	switch {
	case xerror.Is(err, usecase.ErrGrantTypeInvalid):
		return http.StatusBadRequest, OAuth2TokenErrorResponseDTO{
			Error:            "unsupported_grant_type",
			ErrorDescription: err.Error(),
		}

	case xerror.Is(err, usecase.ErrClientInvalid, domain.ErrClientInvalid):
		return http.StatusBadRequest, OAuth2TokenErrorResponseDTO{
			Error:            "invalid_client",
			ErrorDescription: err.Error(),
		}

	case xerror.Is(err, domain.ErrClientUnauthorized):
		return http.StatusUnauthorized, OAuth2TokenErrorResponseDTO{
			Error:            "invalid_client",
			ErrorDescription: err.Error(),
		}

	case xerror.Is(
		err,
		usecase.ErrRefreshTokenInvalid,
		usecase.ErrRefreshTokenStolen,
		usecase.ErrUsernamePasswordInvalid,
	):
		return http.StatusUnauthorized, OAuth2TokenErrorResponseDTO{
			Error:            "invalid_grant",
			ErrorDescription: err.Error(),
		}

	default:
		return 0, OAuth2TokenErrorResponseDTO{}
	}
}
