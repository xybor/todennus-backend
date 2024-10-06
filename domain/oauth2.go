package domain

import (
	"time"

	"github.com/xybor-x/snowflake"
)

type OAuth2TokenMedata struct {
	Id        int64
	Issuer    string
	Audience  string
	Subject   int64
	ExpiresAt int
	NotBefore int

	// Not include in AccessToken
	ExpiresIn int
}

type OAuth2AccessToken struct {
	Metadata OAuth2TokenMedata
}

type OAuth2RefreshToken struct {
	Metadata       OAuth2TokenMedata
	SequenceNumber int
}

type OAuth2IDToken struct {
	Metadata OAuth2TokenMedata
	User     User
}

type OAuth2Domain struct {
	Snowflake *snowflake.Node

	Issuer                 string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	IDTokenExpiration      time.Duration
}

func NewOAuth2Domain(
	snowflake *snowflake.Node,
	issuer string,
	accessTokenExpiration time.Duration,
	refreshTokenExpiration time.Duration,
	idTokenExpiration time.Duration,
) (*OAuth2Domain, error) {
	return &OAuth2Domain{
		Snowflake:              snowflake,
		Issuer:                 issuer,
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
		IDTokenExpiration:      idTokenExpiration,
	}, nil
}

func (domain *OAuth2Domain) CreateAccessToken(aud string, userID int64) (OAuth2AccessToken, error) {
	return OAuth2AccessToken{
		Metadata: domain.createMedata(aud, userID, domain.AccessTokenExpiration),
	}, nil
}

func (domain *OAuth2Domain) CreateRefreshToken(aud string, userID int64) (OAuth2RefreshToken, error) {
	return OAuth2RefreshToken{
		Metadata:       domain.createMedata(aud, userID, domain.RefreshTokenExpiration),
		SequenceNumber: 0,
	}, nil
}

func (domain *OAuth2Domain) NextRefreshToken(current OAuth2RefreshToken) (OAuth2RefreshToken, error) {
	next, err := domain.CreateRefreshToken(current.Metadata.Audience, current.Metadata.Subject)
	if err != nil {
		return OAuth2RefreshToken{}, err
	}

	next.SequenceNumber = current.SequenceNumber + 1
	return next, nil
}

func (domain *OAuth2Domain) CreateIDToken(aud string, user User) (OAuth2IDToken, error) {
	return OAuth2IDToken{
		Metadata: domain.createMedata(aud, user.ID, domain.IDTokenExpiration),
		User:     user,
	}, nil
}

func (domain *OAuth2Domain) createMedata(aud string, sub int64, expiration time.Duration) OAuth2TokenMedata {
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
