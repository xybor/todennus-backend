package abstraction

import (
	"context"

	"github.com/xybor/todennus-backend/usecase/dto"
)

type OAuth2Usecase interface {
	Authorize(ctx context.Context, req *dto.OAuth2AuthorizeRequest) (*dto.OAuth2AuthorizeResponse, error)
	Token(ctx context.Context, req *dto.OAuth2TokenRequest) (*dto.OAuth2TokenResponse, error)
	AuthenticationCallback(ctx context.Context, req *dto.OAuth2AuthenticationCallbackRequest) (*dto.OAuth2AuthenticationCallbackResponse, error)
	SessionUpdate(ctx context.Context, req *dto.OAuth2SessionUpdateRequest) (*dto.OAuth2SessionUpdateResponse, error)
	GetConsent(ctx context.Context, req *dto.OAuth2GetConsentRequest) (*dto.OAuth2GetConsentResponse, error)
	UpdateConsent(ctx context.Context, req *dto.OAuth2UpdateConsentRequest) (*dto.OAUth2UpdateConsentResponse, error)
}
