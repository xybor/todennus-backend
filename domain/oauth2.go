package domain

import (
	"errors"
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/pkg/xrandom"
	"github.com/xybor/todennus-backend/pkg/xstring"
)

const (
	MaximumClientNameLength int = 64
	MinimumClientNameLength int = 3
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
}

type OAuth2RefreshToken struct {
	Metadata       OAuth2TokenMedata
	SequenceNumber int
}

type OAuth2IDToken struct {
	Metadata OAuth2TokenMedata
	User     User
}

type OAuth2AdminToken struct {
	Metadata OAuth2TokenMedata
}

type OAuth2Client struct {
	ID             int64
	OwnerUserID    int64
	Name           string
	HasedSecret    string
	IsConfidential bool
	UpdatedAt      time.Time
}

type OAuth2Domain struct {
	Snowflake *snowflake.Node

	ClientSecretLength int

	Issuer                 string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	IDTokenExpiration      time.Duration
	AdminTokenExpiration   time.Duration
}

func NewOAuth2Domain(
	snowflake *snowflake.Node,
	clientSecretLength int,
	issuer string,
	accessTokenExpiration time.Duration,
	refreshTokenExpiration time.Duration,
	idTokenExpiration time.Duration,
) (*OAuth2Domain, error) {
	return &OAuth2Domain{
		Snowflake:              snowflake,
		ClientSecretLength:     clientSecretLength,
		Issuer:                 issuer,
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
		IDTokenExpiration:      idTokenExpiration,
	}, nil
}

func (domain *OAuth2Domain) CreateClient(ownerID int64, name string, isConfidential bool) (OAuth2Client, string, error) {
	err := domain.validateClientName(name)
	if err != nil {
		return OAuth2Client{}, "", err
	}

	secret := ""
	hashedSecret := []byte{}
	if isConfidential {
		secret = xrandom.RandString(domain.ClientSecretLength)
		hashedSecret, err = HashPassword(secret)
		if err != nil {
			return OAuth2Client{}, "", err
		}
	}

	return OAuth2Client{
		ID:             domain.Snowflake.Generate().Int64(),
		Name:           name,
		OwnerUserID:    ownerID,
		IsConfidential: isConfidential,
		HasedSecret:    string(hashedSecret),
	}, secret, nil
}

func (domain *OAuth2Domain) ValidateClient(
	client OAuth2Client,
	clientID int64,
	clientSecret string,
	confidentialRequirement ConfidentialRequirementType,
) error {
	if client.ID != clientID {
		return errors.New("mismatched client id")
	}

	switch confidentialRequirement {
	case RequireConfidential:
		if !client.IsConfidential {
			return Wrap(ErrClientInvalid, "require a confidential client")
		}

		ok, err := ValidatePassword(client.HasedSecret, clientSecret)
		if err != nil {
			return err
		}

		if !ok {
			return Wrap(ErrClientInvalid, "client secret is invalid")
		}

	case DependOnClientConfidential:
		if client.IsConfidential {
			ok, err := ValidatePassword(client.HasedSecret, clientSecret)
			if err != nil {
				return err
			}

			if !ok {
				return Wrap(ErrClientInvalid, "client secret is invalid")
			}
		}
	}

	return nil
}

func (domain *OAuth2Domain) CreateAccessToken(aud string, user User) (OAuth2AccessToken, error) {
	return OAuth2AccessToken{
		Metadata: domain.createMedata(aud, user.ID, domain.AccessTokenExpiration),
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

	next.Metadata.Id = current.Metadata.Id
	next.SequenceNumber = current.SequenceNumber + 1
	return next, nil
}

func (domain *OAuth2Domain) CreateIDToken(aud string, user User) (OAuth2IDToken, error) {
	return OAuth2IDToken{
		Metadata: domain.createMedata(aud, user.ID, domain.IDTokenExpiration),
		User:     user,
	}, nil
}

func (domain *OAuth2Domain) CreateAdminToken(aud string) (OAuth2AdminToken, error) {
	return OAuth2AdminToken{
		Metadata: domain.createMedata(aud, 0, domain.AdminTokenExpiration),
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

func (domain *OAuth2Domain) validateClientName(clientName string) error {
	if len(clientName) > MaximumClientNameLength {
		return Wrap(ErrClientNameInvalid, "require at most %d characters", MaximumClientNameLength)
	}

	if len(clientName) < MinimumClientNameLength {
		return Wrap(ErrClientNameInvalid, "require at least %d characters", MinimumClientNameLength)
	}

	for _, c := range clientName {
		if !xstring.IsNumber(c) && !xstring.IsLetter(c) && !xstring.IsUnderscore(c) && !xstring.IsSpace(c) {
			return Wrap(ErrClientNameInvalid, "got an invalid character %c", c)
		}
	}

	return nil

}
