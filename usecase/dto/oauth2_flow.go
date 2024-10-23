package dto

import (
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/usecase/dto/resource"
	"github.com/xybor/x/scope"
	"github.com/xybor/x/token"
)

var _ (token.Claims) = (*OAuth2StandardClaims)(nil)

type OAuth2StandardClaims struct {
	ID        string `json:"jti,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	Audience  string `json:"aud,omitempty"`
	Subject   string `json:"sub,omitempty"`
	ExpiresAt int    `json:"exp,omitempty"`
	NotBefore int    `json:"nbf,omitempty"`
}

func OAuth2StandardClaimsFromDomain(claims *domain.OAuth2TokenMedata) *OAuth2StandardClaims {
	return &OAuth2StandardClaims{
		ID:        claims.ID.String(),
		Issuer:    claims.Issuer,
		Audience:  claims.Audience,
		Subject:   claims.Subject.String(),
		ExpiresAt: claims.ExpiresAt,
		NotBefore: claims.NotBefore,
	}
}

func (claims *OAuth2StandardClaims) To() (*domain.OAuth2TokenMedata, error) {
	id, err := snowflake.ParseString(claims.ID)
	if err != nil {
		return nil, err
	}

	sub, err := snowflake.ParseString(claims.Subject)
	if err != nil {
		return nil, err
	}

	return &domain.OAuth2TokenMedata{
		ID:        id,
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

	snowflakeID, err := snowflake.ParseString(claims.ID)
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

func OAuth2AccessTokenFromDomain(token *domain.OAuth2AccessToken) *OAuth2AccessToken {
	return &OAuth2AccessToken{
		OAuth2StandardClaims: OAuth2StandardClaimsFromDomain(token.Metadata),
		Scope:                token.Scope.String(),
	}
}

func (token *OAuth2AccessToken) To() (*domain.OAuth2AccessToken, error) {
	metadata, err := token.OAuth2StandardClaims.To()
	if err != nil {
		return nil, err
	}

	return &domain.OAuth2AccessToken{
		Metadata: metadata,
		Scope:    domain.ScopeEngine.ParseScopes(token.Scope),
	}, nil
}

type OAuth2RefreshToken struct {
	*OAuth2StandardClaims
	SequenceNumber int    `json:"seq"`
	Scope          string `json:"scope"`
}

func OAuth2RefreshTokenFromDomain(token *domain.OAuth2RefreshToken) *OAuth2RefreshToken {
	return &OAuth2RefreshToken{
		OAuth2StandardClaims: OAuth2StandardClaimsFromDomain(token.Metadata),
		SequenceNumber:       token.SequenceNumber,
		Scope:                token.Scope.String(),
	}
}

func (token *OAuth2RefreshToken) To() (*domain.OAuth2RefreshToken, error) {
	metadata, err := token.OAuth2StandardClaims.To()
	if err != nil {
		return nil, err
	}

	return &domain.OAuth2RefreshToken{
		Metadata:       metadata,
		SequenceNumber: token.SequenceNumber,
		Scope:          domain.ScopeEngine.ParseScopes(token.Scope),
	}, nil
}

type OAuth2IDToken struct {
	*OAuth2StandardClaims

	Username    string `json:"username"`
	Displayname string `json:"display_name"`
}

func OAuth2IDTokenFromDomain(token *domain.OAuth2IDToken) *OAuth2IDToken {
	return &OAuth2IDToken{
		OAuth2StandardClaims: OAuth2StandardClaimsFromDomain(token.Metadata),
		Username:             token.User.Username,
		Displayname:          token.User.DisplayName,
	}
}

func (token *OAuth2IDToken) To() (*domain.OAuth2IDToken, error) {
	metadata, err := token.OAuth2StandardClaims.To()
	if err != nil {
		return nil, err
	}

	return &domain.OAuth2IDToken{
		Metadata: metadata,
		User: &domain.User{
			ID:          metadata.Subject,
			Username:    token.Username,
			DisplayName: token.Displayname,
		},
	}, nil
}

type OAuth2TokenRequestDTO struct {
	GrantType string

	ClientID     snowflake.ID
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

type OAuth2AuthorizeRequestDTO struct {
	ResponseType string
	ClientID     snowflake.ID
	RedirectURI  string
	Scope        string
	State        string

	// Only for PKCE
	CodeChallenge       string
	CodeChallengeMethod string
}

type OAuth2AuthorizeResponseDTO struct {
	// Idp
	IdpURL          string
	AuthorizationID string

	// Consent
	NeedConsent bool

	// Authorization Code Flow
	Code string

	// Implicit Flow
	AccessToken string
	TokenType   string
	ExpiresIn   int
}

func NewOAuth2AuthorizeResponseWithCode(code string) *OAuth2AuthorizeResponseDTO {
	return &OAuth2AuthorizeResponseDTO{Code: code}
}

func NewOAuth2AuthorizeResponseRedirectToIdP(url, aid string) *OAuth2AuthorizeResponseDTO {
	return &OAuth2AuthorizeResponseDTO{
		IdpURL:          url,
		AuthorizationID: aid,
	}
}

func NewOAuth2AuthorizeResponseRedirectToConsent(aid string) *OAuth2AuthorizeResponseDTO {
	return &OAuth2AuthorizeResponseDTO{
		NeedConsent:     true,
		AuthorizationID: aid,
	}
}

func NewOAuth2AuthorizeResponseWithToken(token, tokenType string, expiration time.Duration) *OAuth2AuthorizeResponseDTO {
	return &OAuth2AuthorizeResponseDTO{
		AccessToken: token,
		TokenType:   tokenType,
		ExpiresIn:   int(expiration / time.Second),
	}
}

type OAuth2AuthenticationCallbackRequestDTO struct {
	Secret          string
	AuthorizationID string
	Success         bool
	Error           string
	UserID          snowflake.ID
	Username        string
}

type OAuth2AuthenticationCallbackResponseDTO struct {
	AuthenticationID string
}

type OAuth2SessionUpdateRequestDTO struct {
	AuthenticationID string
}

// After updating the session, we must redirect user to Authorization Endpoint
// again. So the response of SessionUpdate is the request of Authorization
// Endpoint.
type OAuth2SessionUpdateResponseDTO OAuth2AuthorizeRequestDTO

func NewOAuth2SessionUpdateResponseDTO(store *domain.OAuth2AuthorizationStore) *OAuth2SessionUpdateResponseDTO {
	return &OAuth2SessionUpdateResponseDTO{
		ResponseType:        store.ResponseType,
		ClientID:            store.ClientID,
		RedirectURI:         store.RedirectURI,
		Scope:               store.Scope.String(),
		State:               store.State,
		CodeChallenge:       store.CodeChallenge,
		CodeChallengeMethod: store.CodeChallengeMethod,
	}
}

type OAuth2GetConsentRequestDTO struct {
	AuthorizationID string
}

type OAuth2GetConsentResponseDTO struct {
	Client *resource.OAuth2Client
	Scopes []string
}

func NewOAuth2GetConsentResponseDTO(client *domain.OAuth2Client, scope scope.Scopes) *OAuth2GetConsentResponseDTO {
	scopeStr := []string{}
	for i := range scope {
		scopeStr = append(scopeStr, scope[i].String())
	}

	return &OAuth2GetConsentResponseDTO{
		Client: resource.NewOAuth2ClientWithoutFilter(client),
		Scopes: scopeStr,
	}
}

type OAuth2UpdateConsentRequestDTO struct {
	AuthorizationID string
	UserScope       string
	Accept          bool
}

// After updating the consent, we must redirect user to Authorization Endpoint
// again. So the response of UpdateConsent is the request of Authorization
// Endpoint.
type OAUth2UpdateConsentResponseDTO OAuth2AuthorizeRequestDTO

func NewOAUth2UpdateConsentResponseDTO(store *domain.OAuth2AuthorizationStore) *OAUth2UpdateConsentResponseDTO {
	return &OAUth2UpdateConsentResponseDTO{
		ResponseType:        store.ResponseType,
		ClientID:            store.ClientID,
		RedirectURI:         store.RedirectURI,
		Scope:               store.Scope.String(),
		State:               store.State,
		CodeChallenge:       store.CodeChallenge,
		CodeChallengeMethod: store.CodeChallengeMethod,
	}
}
