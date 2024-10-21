package abstraction

import (
	"context"

	"github.com/xybor/todennus-backend/usecase/dto"
)

type OAuth2ClientUsecase interface {
	Get(ctx context.Context, req *dto.OAuth2ClientGetRequestDTO) (*dto.OAuth2ClientGetResponseDTO, error)
	Create(ctx context.Context, req *dto.OAuth2ClientCreateRequestDTO) (*dto.OAuth2ClientCreateResponseDTO, error)
	CreateByAdmin(ctx context.Context, req *dto.OAuth2ClientCreateFirstRequestDTO) (*dto.OAuth2ClientCreateByAdminResponseDTO, error)
}
