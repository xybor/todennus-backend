package domain

import (
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/pkg/scope"
)

type ConfidentialRequirementType int

const (
	RequireConfidential ConfidentialRequirementType = iota
	NotRequireConfidential
	DependOnClientConfidential
)

type OAuth2TokenMedata struct {
	Id        int64
	Issuer    string
	Audience  string
	Subject   int64
	ExpiresAt int
	NotBefore int

	// Not included in AccessToken
	ExpiresIn int
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

type OAuth2AdminToken struct {
	Metadata OAuth2TokenMedata
}

type OAuth2FlowDomain struct {
	Snowflake *snowflake.Node

	Issuer                 string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	IDTokenExpiration      time.Duration
	AdminTokenExpiration   time.Duration
}

func NewOAuth2FlowDomain(
	snowflake *snowflake.Node,
	issuer string,
	accessTokenExpiration time.Duration,
	refreshTokenExpiration time.Duration,
	idTokenExpiration time.Duration,
) (*OAuth2FlowDomain, error) {
	return &OAuth2FlowDomain{
		Snowflake:              snowflake,
		Issuer:                 issuer,
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
		IDTokenExpiration:      idTokenExpiration,
	}, nil
}

func (domain *OAuth2FlowDomain) CreateAccessToken(aud string, scope scope.Scopes, user User) (OAuth2AccessToken, error) {
	return OAuth2AccessToken{
		Metadata: domain.createMedata(aud, user.ID, domain.AccessTokenExpiration),
		Scope:    scope,
	}, nil
}

func (domain *OAuth2FlowDomain) CreateRefreshToken(aud string, scope scope.Scopes, userID int64) (OAuth2RefreshToken, error) {
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

	next.Metadata.Id = current.Metadata.Id
	next.SequenceNumber = current.SequenceNumber + 1
	return next, nil
}

func (domain *OAuth2FlowDomain) CreateIDToken(aud string, user User) (OAuth2IDToken, error) {
	return OAuth2IDToken{
		Metadata: domain.createMedata(aud, user.ID, domain.IDTokenExpiration),
		User:     user,
	}, nil
}

func (domain *OAuth2FlowDomain) CreateAdminToken(aud string) (OAuth2AdminToken, error) {
	return OAuth2AdminToken{
		Metadata: domain.createMedata(aud, 0, domain.AdminTokenExpiration),
	}, nil
}

func (domain *OAuth2FlowDomain) createMedata(aud string, sub int64, expiration time.Duration) OAuth2TokenMedata {
	id := domain.Snowflake.Generate()

	return OAuth2TokenMedata{
		Id:        id.Int64(),
		Issuer:    domain.Issuer,
		Audience:  aud,
		Subject:   sub,
		ExpiresAt: int(time.UnixMilli(id.Time()).Add(expiration).Unix()),
		NotBefore: int(time.UnixMilli(id.Time()).Unix()),

		ExpiresIn: int(expiration / time.Second),
	}
}
