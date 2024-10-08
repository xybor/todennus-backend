package abstraction

import "github.com/xybor/todennus-backend/domain"

type UserDomain interface {
	Create(username, password string) (domain.User, error)
	Validate(hashedPassword, password string) (bool, error)
}

type OAuth2Domain interface {
	CreateClient(ownerID int64, name string, isConfidential bool) (domain.OAuth2Client, string, error)
	ValidateClient(
		client domain.OAuth2Client,
		clientID int64,
		clientSecret string,
		confidentialRequirement domain.ConfidentialRequirementType,
	) error

	CreateAccessToken(aud string, user domain.User) (domain.OAuth2AccessToken, error)
	CreateRefreshToken(aud string, userID int64) (domain.OAuth2RefreshToken, error)
	NextRefreshToken(current domain.OAuth2RefreshToken) (domain.OAuth2RefreshToken, error)
	CreateIDToken(aud string, user domain.User) (domain.OAuth2IDToken, error)
}
