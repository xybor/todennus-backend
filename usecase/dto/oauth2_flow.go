package dto

import (
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/pkg/token"
	"github.com/xybor/todennus-backend/pkg/xstring"
)

var _ (token.Claims) = (*OAuth2StandardClaims)(nil)

type OAuth2StandardClaims struct {
	Id        string `json:"jti,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	Audience  string `json:"aud,omitempty"`
	Subject   string `json:"sub,omitempty"`
	ExpiresAt int    `json:"exp,omitempty"`
	NotBefore int    `json:"nbf,omitempty"`
}

func OAuth2StandardClaimsFromDomain(claims domain.OAuth2TokenMedata) *OAuth2StandardClaims {
	return &OAuth2StandardClaims{
		Id:        xstring.FormatID(claims.Id),
		Issuer:    claims.Issuer,
		Audience:  claims.Audience,
		Subject:   xstring.FormatID(claims.Subject),
		ExpiresAt: claims.ExpiresAt,
		NotBefore: claims.NotBefore,
	}
}

func (claims *OAuth2StandardClaims) To() (domain.OAuth2TokenMedata, error) {
	id, err := xstring.ParseID(claims.Id)
	if err != nil {
		return domain.OAuth2TokenMedata{}, err
	}

	sub, err := xstring.ParseID(claims.Subject)
	if err != nil {
		return domain.OAuth2TokenMedata{}, err
	}

	return domain.OAuth2TokenMedata{
		Id:        id,
		Issuer:    claims.Issuer,
		Audience:  claims.Audience,
		Subject:   sub,
		ExpiresAt: claims.ExpiresAt,
		NotBefore: claims.NotBefore,
	}, nil
}

func (claims *OAuth2StandardClaims) Valid() error {
	now := time.Now()
	if claims.ExpiresAt != 0 && time.Unix(int64(claims.ExpiresAt), 0).Before(now) {
		return token.ErrTokenExpired
	}

	if claims.NotBefore != 0 && time.Unix(int64(claims.NotBefore), 0).After(now) {
		return token.ErrTokenNotYetValid
	}

	snowflakeID, err := snowflake.ParseString(claims.Id)
	if err != nil {
		return token.ErrTokenInvalidFormat
	}

	createdAt := time.UnixMilli(snowflakeID.Time())
	if createdAt.After(now) {
		return token.ErrTokenNotYetValid
	}

	return nil
}

type OAuth2AccessToken struct {
	*OAuth2StandardClaims
	Scope string `json:"scope"`
}

func OAuth2AccessTokenFromDomain(token domain.OAuth2AccessToken) OAuth2AccessToken {
	return OAuth2AccessToken{
		OAuth2StandardClaims: OAuth2StandardClaimsFromDomain(token.Metadata),
		Scope:                token.Scope.String(),
	}
}

func (token *OAuth2AccessToken) To() (domain.OAuth2AccessToken, error) {
	scope, err := domain.ScopeEngine.ParseScopes(token.Scope)
	if err != nil {
		return domain.OAuth2AccessToken{}, err
	}

	metadata, err := token.OAuth2StandardClaims.To()
	if err != nil {
		return domain.OAuth2AccessToken{}, err
	}

	return domain.OAuth2AccessToken{
		Metadata: metadata,
		Scope:    scope,
	}, nil
}

type OAuth2RefreshToken struct {
	*OAuth2StandardClaims
	SequenceNumber int    `json:"seq"`
	Scope          string `json:"scope"`
}

func OAuth2RefreshTokenFromDomain(token domain.OAuth2RefreshToken) OAuth2RefreshToken {
	return OAuth2RefreshToken{
		OAuth2StandardClaims: OAuth2StandardClaimsFromDomain(token.Metadata),
		SequenceNumber:       token.SequenceNumber,
		Scope:                token.Scope.String(),
	}
}

func (token *OAuth2RefreshToken) To() (domain.OAuth2RefreshToken, error) {
	scope, err := domain.ScopeEngine.ParseScopes(token.Scope)
	if err != nil {
		return domain.OAuth2RefreshToken{}, err
	}

	metadata, err := token.OAuth2StandardClaims.To()
	if err != nil {
		return domain.OAuth2RefreshToken{}, err
	}

	return domain.OAuth2RefreshToken{
		Metadata:       metadata,
		SequenceNumber: token.SequenceNumber,
		Scope:          scope,
	}, nil
}

type OAuth2IDToken struct {
	*OAuth2StandardClaims

	Username    string `json:"username"`
	Displayname string `json:"display_name"`
}

func OAuth2IDTokenFromDomain(token domain.OAuth2IDToken) OAuth2IDToken {
	return OAuth2IDToken{
		OAuth2StandardClaims: OAuth2StandardClaimsFromDomain(token.Metadata),
		Username:             token.User.Username,
		Displayname:          token.User.DisplayName,
	}
}

func (token *OAuth2IDToken) To() (domain.OAuth2IDToken, error) {
	metadata, err := token.OAuth2StandardClaims.To()
	if err != nil {
		return domain.OAuth2IDToken{}, err
	}

	return domain.OAuth2IDToken{
		Metadata: metadata,
		User: domain.User{
			ID:          metadata.Subject,
			Username:    token.Username,
			DisplayName: token.Displayname,
		},
	}, nil
}

type OAuth2TokenRequestDTO struct {
	GrantType string

	ClientID     int64
	ClientSecret string

	// Authorization Code Flow
	Code         string
	RedirectURI  string
	CodeVerifier string // with PKCE

	// Resource Owner Password Credentials Flow
	Username string
	Password string
	Scope    string

	// Refresh Token Flow
	RefreshToken string
}

type OAuth2TokenResponseDTO struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	RefreshToken string
	Scope        string
}
