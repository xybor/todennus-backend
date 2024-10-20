package domain

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/x/scope"
	"github.com/xybor/x/xcrypto"
)

type SessionState int

const (
	SessionStateUnauthenticated = iota
	SessionStateAuthenticated
	SessionStateFailedAuthentication
)

const (
	CodeChallengeMethodPlain = "plain"
	CodeChallengeMethodS256  = "S256"
)

type Session struct {
	State     SessionState
	UserID    snowflake.ID
	ExpiresAt time.Time
}

type OAuth2AuthorizationCode struct {
	Code                string
	UserID              snowflake.ID
	ClientID            snowflake.ID
	Scope               scope.Scopes
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
}

type OAuth2AuthorizationStore struct {
	ID                  string
	HasAuthenticated    bool
	ResponseType        string
	ClientID            snowflake.ID
	RedirectURI         string
	Scope               scope.Scopes
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
}

type OAuth2AuthenticationResult struct {
	ID              string
	AuthorizationID string
	Ok              bool
	Error           string
	UserID          snowflake.ID
	Username        string
	ExpiresAt       time.Time
}

type OAuth2TokenMedata struct {
	ID        snowflake.ID
	Issuer    string
	Audience  string
	Subject   snowflake.ID
	ExpiresAt int
	NotBefore int
}

type OAuth2AccessToken struct {
	Metadata OAuth2TokenMedata
	Scope    scope.Scopes
}

type OAuth2RefreshToken struct {
	Metadata       OAuth2TokenMedata
	SequenceNumber int
	Scope          scope.Scopes
}

type OAuth2IDToken struct {
	Metadata OAuth2TokenMedata
	User     User
}

type OAuth2FlowDomain struct {
	Snowflake *snowflake.Node
	Issuer    string

	AuthorizationCodeFlowExpiration time.Duration
	LoginResultExpiration           time.Duration
	LoginUpdateExpiration           time.Duration
	SessionExpiration               time.Duration

	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	IDTokenExpiration      time.Duration
	AdminTokenExpiration   time.Duration
}

func NewOAuth2FlowDomain(
	snowflake *snowflake.Node,
	issuer string,
	authorizationCodeFlowExpiration time.Duration,
	loginResultExpiration time.Duration,
	loginUpdateExpiration time.Duration,
	sessionExpiration time.Duration,
	accessTokenExpiration time.Duration,
	refreshTokenExpiration time.Duration,
	idTokenExpiration time.Duration,
) (*OAuth2FlowDomain, error) {
	return &OAuth2FlowDomain{
		Snowflake: snowflake,
		Issuer:    issuer,

		AuthorizationCodeFlowExpiration: accessTokenExpiration,
		LoginResultExpiration:           loginResultExpiration,
		LoginUpdateExpiration:           loginUpdateExpiration,
		SessionExpiration:               sessionExpiration,

		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
		IDTokenExpiration:      idTokenExpiration,
	}, nil
}

func (domain *OAuth2FlowDomain) CreateAuthorizationCode(
	userID, clientID snowflake.ID,
	scope scope.Scopes,
	codeChallenge, codeChallengeMethod string,
) OAuth2AuthorizationCode {
	return OAuth2AuthorizationCode{
		Code:                xcrypto.RandString(32),
		Scope:               scope,
		UserID:              userID,
		ClientID:            clientID,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		ExpiresAt:           time.Now().Add(domain.AuthorizationCodeFlowExpiration),
	}
}

func (domain *OAuth2FlowDomain) CreateAuthorizationStore(
	respType string,
	clientID snowflake.ID,
	scope scope.Scopes,
	redirectURI, state, codeChallenge, codeChallengeMethod string,
) OAuth2AuthorizationStore {
	return OAuth2AuthorizationStore{
		ID:                  xcrypto.RandString(32),
		ResponseType:        respType,
		HasAuthenticated:    false,
		Scope:               scope,
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		ExpiresAt:           time.Now().Add(domain.LoginResultExpiration),
	}
}

func (domain *OAuth2FlowDomain) CreateAuthenticationResultSuccess(authID string, userID snowflake.ID, username string) OAuth2AuthenticationResult {
	return OAuth2AuthenticationResult{
		ID:              xcrypto.RandString(32),
		Ok:              true,
		AuthorizationID: authID,
		UserID:          userID,
		Username:        username,
		ExpiresAt:       time.Now().Add(domain.LoginUpdateExpiration),
	}
}

func (domain *OAuth2FlowDomain) CreateAuthenticationResultFailure(authID string, err string) OAuth2AuthenticationResult {
	return OAuth2AuthenticationResult{
		ID:              xcrypto.RandString(32),
		Ok:              false,
		AuthorizationID: authID,
		Error:           err,
		ExpiresAt:       time.Now().Add(domain.LoginUpdateExpiration),
	}
}

func (domain *OAuth2FlowDomain) CreateAccessToken(aud string, scope scope.Scopes, user User) (OAuth2AccessToken, error) {
	return OAuth2AccessToken{
		Metadata: domain.createMedata(aud, user.ID, domain.AccessTokenExpiration),
		Scope:    scope,
	}, nil
}

func (domain *OAuth2FlowDomain) CreateRefreshToken(aud string, scope scope.Scopes, userID snowflake.ID) (OAuth2RefreshToken, error) {
	return OAuth2RefreshToken{
		Metadata:       domain.createMedata(aud, userID, domain.RefreshTokenExpiration),
		SequenceNumber: 0,
		Scope:          scope,
	}, nil
}

func (domain *OAuth2FlowDomain) NextRefreshToken(current OAuth2RefreshToken) (OAuth2RefreshToken, error) {
	next, err := domain.CreateRefreshToken(current.Metadata.Audience, current.Scope, current.Metadata.Subject)
	if err != nil {
		return OAuth2RefreshToken{}, err
	}

	next.Metadata.ID = current.Metadata.ID
	next.SequenceNumber = current.SequenceNumber + 1
	return next, nil
}

func (domain *OAuth2FlowDomain) CreateIDToken(aud string, user User) (OAuth2IDToken, error) {
	return OAuth2IDToken{
		Metadata: domain.createMedata(aud, user.ID, domain.IDTokenExpiration),
		User:     user,
	}, nil
}

func (domain *OAuth2FlowDomain) ValidateCodeChallenge(verifier, challenge, method string) bool {
	switch method {
	case CodeChallengeMethodPlain:
		return verifier == challenge
	default: // CodeChallengeMethodS256
		hash := sha256.Sum256([]byte(verifier))
		encoded := base64.RawURLEncoding.EncodeToString(hash[:])
		return encoded == challenge
	}
}

func (domain *OAuth2FlowDomain) NewSession(userID snowflake.ID) Session {
	return Session{
		State:     SessionStateAuthenticated,
		UserID:    userID,
		ExpiresAt: time.Now().Add(domain.SessionExpiration),
	}
}

func (domain *OAuth2FlowDomain) InvalidateSession(state SessionState) Session {
	if state != SessionStateFailedAuthentication && state != SessionStateUnauthenticated {
		panic("invalid call")
	}

	return Session{State: state, ExpiresAt: time.Now().Add(domain.SessionExpiration)}
}

func (domain *OAuth2FlowDomain) createMedata(aud string, sub snowflake.ID, expiration time.Duration) OAuth2TokenMedata {
	id := domain.Snowflake.Generate()

	return OAuth2TokenMedata{
		ID:        id,
		Issuer:    domain.Issuer,
		Audience:  aud,
		Subject:   sub,
		ExpiresAt: int(time.UnixMilli(id.Time()).Add(expiration).Unix()),
		NotBefore: int(time.UnixMilli(id.Time()).Unix()),
	}
}
