package dto

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/todennus-backend/usecase/dto"
	"github.com/xybor/x/logging"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xerror"
	"github.com/xybor/x/xhttp"
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
	Scope    string `form:"scope"`

	// Refresh Token Flow
	RefreshToken string `form:"refresh_token"`
}

func (req OAuth2TokenRequestDTO) To() dto.OAuth2TokenRequestDTO {
	return dto.OAuth2TokenRequestDTO{
		GrantType: req.GrantType,

		ClientID:     snowflake.ID(req.ClientID),
		ClientSecret: req.ClientSecret,

		Code:         req.Code,
		RedirectURI:  req.RedirectURI,
		CodeVerifier: req.CodeVerifier,

		Username: req.Username,
		Password: req.Password,
		Scope:    req.Scope,

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
	case xerror.Is(err, usecase.ErrTokenInvalidGrantType):
		return http.StatusBadRequest, OAuth2TokenErrorResponseDTO{
			Error:            "unsupported_grant_type",
			ErrorDescription: xerror.MessageOf(err),
		}

	case xerror.Is(err, usecase.ErrClientInvalid, domain.ErrClientInvalid, domain.ErrClientUnauthorized):
		return http.StatusBadRequest, OAuth2TokenErrorResponseDTO{
			Error:            "invalid_client",
			ErrorDescription: xerror.MessageOf(err),
		}

	case xerror.Is(err, usecase.ErrScopeInvalid):
		return http.StatusBadRequest, OAuth2TokenErrorResponseDTO{
			Error:            "invalid_scope",
			ErrorDescription: xerror.MessageOf(err),
		}

	case xerror.Is(err, usecase.ErrTokenInvalidGrant):
		return http.StatusUnauthorized, OAuth2TokenErrorResponseDTO{
			Error:            "invalid_grant",
			ErrorDescription: xerror.MessageOf(err),
		}

	default:
		return 0, OAuth2TokenErrorResponseDTO{}
	}
}

type OAuth2AuthorizeRequestDTO struct {
	ResponseType string `query:"response_type"`
	ClientID     int64  `query:"client_id"`
	RedirectURI  string `query:"redirect_uri"`
	Scope        string `query:"scope"`
	State        string `query:"state"`

	// For PKCE
	CodeChallenge       string `query:"code_challenge"`
	CodeChallengeMethod string `query:"code_challenge_method"`
}

func (req OAuth2AuthorizeRequestDTO) To() dto.OAuth2AuthorizeRequestDTO {
	return dto.OAuth2AuthorizeRequestDTO{
		ResponseType:        req.ResponseType,
		ClientID:            snowflake.ID(req.ClientID),
		RedirectURI:         req.RedirectURI,
		Scope:               req.Scope,
		State:               req.State,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
	}
}

func NewOAuth2AuthorizeRedirectURI(req OAuth2AuthorizeRequestDTO, resp dto.OAuth2AuthorizeResponseDTO) (string, error) {
	if resp.IdpURL != "" {
		u, err := url.Parse(resp.IdpURL)
		if err != nil {
			return "", err
		}

		q := u.Query()
		q.Set("authorization_id", resp.AuthorizationID)
		u.RawQuery = q.Encode()

		return u.String(), nil
	}

	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		return "", err
	}

	q := u.Query()

	if req.State != "" {
		q.Set("state", req.State)
	}

	if resp.Code != "" {
		q.Set("code", resp.Code)
	} else {
		if resp.AccessToken != "" || resp.ExpiresIn == 0 || resp.TokenType == "" {
			return "", errors.New("expected access_token, token_type, and expires_in if resp type is token")
		}

		q.Set("access_token", resp.AccessToken)
		q.Set("token_type", resp.TokenType)
		q.Set("expires_in", strconv.FormatInt(int64(resp.ExpiresIn), 10))
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func NewOAuth2AuthorizeRedirectURIWithError(ctx context.Context, req OAuth2AuthorizeRequestDTO, err error) (string, error) {
	u, uerr := xhttp.ParseURL(req.RedirectURI)
	if uerr != nil {
		return "", err
	}

	q := u.Query()

	if req.State != "" {
		q.Set("state", req.State)
	}

	switch {
	case xerror.Is(err, domain.ErrClientInvalid, usecase.ErrClientInvalid):
		q.Set("error", "invalid_client")
		q.Set("error_description", xerror.MessageOf(err))
	case xerror.Is(err, usecase.ErrAuthorizationAccessDenied):
		q.Set("error", "access_denied")
		q.Set("error_description", xerror.MessageOf(err))
	case xerror.Is(err,
		usecase.ErrRequestInvalid,
		usecase.ErrAuthorizationResponseTypeInvalid,
	):
		q.Set("error", "invalid_request")
		q.Set("error_description", xerror.MessageOf(err))
	default:
		q.Set("error", "server_error")
		q.Set("error_description", "the authorization server encountered an unexpected condition")
		logging.LogError(xcontext.Logger(ctx), err)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

type OAuth2AuthenticationCallbackRequestDTO struct {
	IdPSecret       string `json:"idp_secret"`
	AuthorizationID string `json:"authorization_id"`
	Success         bool   `json:"success"`
	UserID          string `json:"user_id"`
	Username        string `json:"username"`
	Error           string `json:"error"`
}

func (req OAuth2AuthenticationCallbackRequestDTO) To() (dto.OAuth2AuthenticationCallbackRequestDTO, error) {
	uid, err := snowflake.ParseString(req.UserID)
	if err != nil {
		return dto.OAuth2AuthenticationCallbackRequestDTO{}, err
	}

	return dto.OAuth2AuthenticationCallbackRequestDTO{
		Secret:          req.IdPSecret,
		Success:         req.Success,
		AuthorizationID: req.AuthorizationID,
		UserID:          uid,
		Username:        req.Username,
		Error:           req.Error,
	}, nil
}

type OAuth2AuthenticationCallbackResponseDTO struct {
	AuthenticationID string `json:"authentication_id"`
}

func NewOAuth2AuthenticationCallbackResponseDTO(resp dto.OAuth2AuthenticationCallbackResponseDTO) OAuth2AuthenticationCallbackResponseDTO {
	return OAuth2AuthenticationCallbackResponseDTO{
		AuthenticationID: resp.AuthenticationID,
	}
}

type OAuth2SessionUpdateRequestDTO struct {
	AuthenticationID string `query:"authentication_id"`
}

func (req OAuth2SessionUpdateRequestDTO) To() dto.OAuth2SessionUpdateRequestDTO {
	return dto.OAuth2SessionUpdateRequestDTO{
		AuthenticationID: req.AuthenticationID,
	}
}

func NewOAuth2SessionUpdateRedirectURI(resp dto.OAuth2SessionUpdateResponseDTO) string {
	q := url.Values{}
	q.Set("response_type", resp.ResponseType)
	q.Set("client_id", resp.ClientID.String())
	q.Set("redirect_uri", resp.RedirectURI)
	q.Set("scope", resp.Scope)

	if resp.State != "" {
		q.Set("state", resp.State)
	}

	if resp.CodeChallenge != "" {
		q.Set("code_challenge", resp.CodeChallenge)
	}

	if resp.CodeChallengeMethod != "" {
		q.Set("code_challenge_method", resp.CodeChallengeMethod)
	}

	return fmt.Sprintf("/oauth2/authorize?%s", q.Encode())
}
