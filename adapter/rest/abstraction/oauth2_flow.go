package abstraction

import (
	"context"

	"github.com/xybor/todennus-backend/usecase/dto"
)

type OAuth2Usecase interface {
	Authorize(ctx context.Context, req *dto.OAuth2AuthorizeRequestDTO) (*dto.OAuth2AuthorizeResponseDTO, error)
	Token(ctx context.Context, req *dto.OAuth2TokenRequestDTO) (*dto.OAuth2TokenResponseDTO, error)
	AuthenticationCallback(ctx context.Context, req *dto.OAuth2AuthenticationCallbackRequestDTO) (*dto.OAuth2AuthenticationCallbackResponseDTO, error)
	SessionUpdate(ctx context.Context, req *dto.OAuth2SessionUpdateRequestDTO) (*dto.OAuth2SessionUpdateResponseDTO, error)
	GetConsent(ctx context.Context, req *dto.OAuth2GetConsentRequestDTO) (*dto.OAuth2GetConsentResponseDTO, error)
	UpdateConsent(ctx context.Context, req *dto.OAuth2UpdateConsentRequestDTO) (*dto.OAUth2UpdateConsentResponseDTO, error)
}
