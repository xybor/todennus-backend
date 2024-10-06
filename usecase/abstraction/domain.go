package abstraction

import "github.com/xybor/todennus-backend/domain"

type UserDomain interface {
	Create(username, password string) (domain.User, error)
	Validate(hashedPassword, password string) (bool, error)
}

type OAuth2Domain interface {
	CreateAccessToken(aud string, userID int64) (domain.OAuth2AccessToken, error)
	CreateRefreshToken(aud string, userID int64) (domain.OAuth2RefreshToken, error)
	NextRefreshToken(current domain.OAuth2RefreshToken) (domain.OAuth2RefreshToken, error)
	CreateIDToken(aud string, user domain.User) (domain.OAuth2IDToken, error)
}
