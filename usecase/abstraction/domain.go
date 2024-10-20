package abstraction

import (
	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/x/scope"
)

type UserDomain interface {
	Create(username, password string) (domain.User, error)
	Validate(hashedPassword, password string) (bool, error)
}

type OAuth2FlowDomain interface {
	CreateAuthorizationCode(
		userID, clientID snowflake.ID,
		scope scope.Scopes,
		codeChallenge, codeChallengeMethod string,
	) domain.OAuth2AuthorizationCode
	CreateAuthorizationStore(
		respType string,
		clientID snowflake.ID,
		scope scope.Scopes,
		redirectURI, state, codeChallenge, codeChallengeMethod string,
	) domain.OAuth2AuthorizationStore
	CreateAuthenticationResultSuccess(authID string, userID snowflake.ID, username string) domain.OAuth2AuthenticationResult
	CreateAuthenticationResultFailure(authID string, err string) domain.OAuth2AuthenticationResult

	CreateAccessToken(aud string, scope scope.Scopes, user domain.User) domain.OAuth2AccessToken
	CreateRefreshToken(aud string, scope scope.Scopes, userID snowflake.ID) domain.OAuth2RefreshToken
	NextRefreshToken(current domain.OAuth2RefreshToken) domain.OAuth2RefreshToken
	CreateIDToken(aud string, user domain.User) domain.OAuth2IDToken

	ValidateCodeChallenge(verifier, challenge, method string) bool

	NewSession(userID snowflake.ID) domain.Session
	InvalidateSession(state domain.SessionState) domain.Session
}

type OAuth2ClientDomain interface {
	CreateClient(ownerID snowflake.ID, name string, isConfidential bool) (domain.OAuth2Client, string, error)
	ValidateClient(
		client domain.OAuth2Client,
		clientID snowflake.ID,
		clientSecret string,
		confidentialRequirement domain.ConfidentialRequirementType,
	) error
}
