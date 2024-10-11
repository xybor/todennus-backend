package domain

import (
	"errors"
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/pkg/scope"
	"github.com/xybor/todennus-backend/pkg/xrandom"
	"github.com/xybor/todennus-backend/pkg/xstring"
)

const (
	MaximumClientNameLength int = 64
	MinimumClientNameLength int = 3
)

type OAuth2Client struct {
	ID             int64
	OwnerUserID    int64
	Name           string
	HashedSecret   string
	IsConfidential bool
	AllowedScope   scope.Scopes
	UpdatedAt      time.Time
}

type OAuth2ClientDomain struct {
	Snowflake          *snowflake.Node
	ClientSecretLength int
}

func NewOAuth2ClientDomain(
	snowflake *snowflake.Node,
	clientSecretLength int,
) (*OAuth2ClientDomain, error) {
	return &OAuth2ClientDomain{
		Snowflake:          snowflake,
		ClientSecretLength: clientSecretLength,
	}, nil
}

func (domain *OAuth2ClientDomain) CreateClient(ownerID int64, name string, isConfidential bool) (OAuth2Client, string, error) {
	err := domain.validateClientName(name)
	if err != nil {
		return OAuth2Client{}, "", err
	}

	secret := ""
	allowedScope := ScopeEngine.New(Actions.Read, Resources).AsScopes()
	hashedSecret := []byte{}
	if isConfidential {
		secret = xrandom.RandString(domain.ClientSecretLength)
		hashedSecret, err = HashPassword(secret)
		if err != nil {
			return OAuth2Client{}, "", err
		}

		allowedScope = ScopeEngine.New(Actions, Resources).AsScopes()
	}

	return OAuth2Client{
		ID:             domain.Snowflake.Generate().Int64(),
		Name:           name,
		OwnerUserID:    ownerID,
		IsConfidential: isConfidential,
		AllowedScope:   allowedScope,
		HashedSecret:   string(hashedSecret),
	}, secret, nil
}

func (domain *OAuth2ClientDomain) ValidateClient(
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

		ok, err := ValidatePassword(client.HashedSecret, clientSecret)
		if err != nil {
			return err
		}

		if !ok {
			return Wrap(ErrClientInvalid, "client secret is invalid")
		}

	case DependOnClientConfidential:
		if client.IsConfidential {
			ok, err := ValidatePassword(client.HashedSecret, clientSecret)
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

func (domain *OAuth2ClientDomain) validateClientName(clientName string) error {
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
