package dto

import (
	"strconv"
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/pkg/token"
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
		Id:        strconv.FormatInt(claims.Id, 10),
		Issuer:    claims.Issuer,
		Audience:  claims.Audience,
		Subject:   strconv.FormatInt(claims.Subject, 10),
		ExpiresAt: claims.ExpiresAt,
		NotBefore: claims.NotBefore,
	}
}

func (claims *OAuth2StandardClaims) To() domain.OAuth2TokenMedata {
	// We signed the claims, so we ensure that this claims is always in correct
	// form.
	id, err := strconv.ParseInt(claims.Id, 10, 64)
	if err != nil {
		panic(err)
	}

	sub, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		panic(err)
	}

	return domain.OAuth2TokenMedata{
		Id:        id,
		Issuer:    claims.Issuer,
		Audience:  claims.Audience,
		Subject:   sub,
		ExpiresAt: claims.ExpiresAt,
		NotBefore: claims.NotBefore,
	}
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
}

func OAuth2AccessTokenFromDomain(token domain.OAuth2AccessToken) OAuth2AccessToken {
	return OAuth2AccessToken{
		OAuth2StandardClaims: OAuth2StandardClaimsFromDomain(token.Metadata),
	}
}

func (token *OAuth2AccessToken) To() domain.OAuth2AccessToken {
	return domain.OAuth2AccessToken{
		Metadata: token.OAuth2StandardClaims.To(),
	}
}

type OAuth2RefreshToken struct {
	*OAuth2StandardClaims
	SequenceNumber int `json:"seq"`
}

func OAuth2RefreshTokenFromDomain(token domain.OAuth2RefreshToken) OAuth2RefreshToken {
	return OAuth2RefreshToken{
		OAuth2StandardClaims: OAuth2StandardClaimsFromDomain(token.Metadata),
		SequenceNumber:       token.SequenceNumber,
	}
}

func (token *OAuth2RefreshToken) To() domain.OAuth2RefreshToken {
	return domain.OAuth2RefreshToken{
		Metadata:       token.OAuth2StandardClaims.To(),
		SequenceNumber: token.SequenceNumber,
	}
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

func (token *OAuth2IDToken) To() domain.OAuth2IDToken {
	metadata := token.OAuth2StandardClaims.To()

	return domain.OAuth2IDToken{
		Metadata: metadata,
		User: domain.User{
			ID:          metadata.Subject,
			Username:    token.Username,
			DisplayName: token.Displayname,
		},
	}
}

type OAuth2TokenRequest struct {
	GrantType string

	ClientID     string
	ClientSecret string

	// Authorization Code Flow
	Code         string
	RedirectURI  string
	CodeVerifier string // with PKCE

	// Resource Owner Password Credentials Flow
	Username string
	Password string

	// Refresh Token Flow
	RefreshToken string
}

type OAuth2TokenResponse struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	RefreshToken string
	Scope        string
}
