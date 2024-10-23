package domain

import (
	"fmt"
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/x/scope"
)

type OAuth2ConsentResult struct {
	UserID    snowflake.ID
	ClientID  snowflake.ID
	Accepted  bool
	Scope     scope.Scopes
	ExpiresAt time.Time
}

type OAuth2Consent struct {
	UserID    snowflake.ID
	ClientID  snowflake.ID
	Scope     scope.Scopes
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type OAuth2ConsentDomain struct {
	FailureConsentExpiration time.Duration
	ConsentExpiration        time.Duration
}

func NewOAuth2ConsentDomain(
	failureExpiration, consentExpiration time.Duration,
) *OAuth2ConsentDomain {
	return &OAuth2ConsentDomain{
		FailureConsentExpiration: failureExpiration,
		ConsentExpiration:        consentExpiration,
	}
}

func (domain *OAuth2ConsentDomain) CreateConsentDeniedResult(userID, clientID snowflake.ID) *OAuth2ConsentResult {
	return &OAuth2ConsentResult{
		Accepted:  false,
		UserID:    userID,
		ClientID:  clientID,
		ExpiresAt: time.Now().Add(domain.FailureConsentExpiration),
	}
}

func (domain *OAuth2ConsentDomain) CreateConsentAcceptedResult(userID, clientID snowflake.ID, scope scope.Scopes) *OAuth2ConsentResult {
	return &OAuth2ConsentResult{
		Accepted:  true,
		UserID:    userID,
		ClientID:  clientID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(domain.FailureConsentExpiration),
	}
}

func (domain *OAuth2ConsentDomain) CreateConsent(
	userID, clientID snowflake.ID,
	scope scope.Scopes,
) *OAuth2Consent {
	return &OAuth2Consent{
		UserID:    userID,
		ClientID:  clientID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(domain.ConsentExpiration),
	}
}

func (domain *OAuth2ConsentDomain) ValidateConsent(consent *OAuth2Consent, requestedScope scope.Scopes) error {
	if consent.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("%w%s", ErrKnown, "consent expired")
	}

	if !requestedScope.LessThanOrEqual(consent.Scope) {
		return fmt.Errorf("%w%s", ErrKnown, "requested scope exceeds the consented scope")
	}

	return nil
}
