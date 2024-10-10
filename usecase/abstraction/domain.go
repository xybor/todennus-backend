package abstraction

import (
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/pkg/scope"
)

type UserDomain interface {
	Create(username, password string) (domain.User, error)
	Validate(hashedPassword, password string) (bool, error)
}

type OAuth2FlowDomain interface {
	CreateAccessToken(aud string, scope scope.Scopes, user domain.User) (domain.OAuth2AccessToken, error)
	CreateRefreshToken(aud string, scope scope.Scopes, userID int64) (domain.OAuth2RefreshToken, error)
	NextRefreshToken(current domain.OAuth2RefreshToken) (domain.OAuth2RefreshToken, error)
	CreateIDToken(aud string, user domain.User) (domain.OAuth2IDToken, error)
}

type OAuth2ClientDomain interface {
	CreateClient(ownerID int64, name string, isConfidential bool) (domain.OAuth2Client, string, error)
	ValidateClient(
		client domain.OAuth2Client,
		clientID int64,
		clientSecret string,
		confidentialRequirement domain.ConfidentialRequirementType,
	) error
}
