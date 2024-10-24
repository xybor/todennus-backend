package abstraction

import (
	"context"

	"github.com/xybor/todennus-backend/usecase/dto"
)

type OAuth2ClientUsecase interface {
	Get(ctx context.Context, req *dto.OAuth2ClientGetRequest) (*dto.OAuth2ClientGetResponse, error)
	Create(ctx context.Context, req *dto.OAuth2ClientCreateRequest) (*dto.OAuth2ClientCreateResponse, error)
	CreateByAdmin(ctx context.Context, req *dto.OAuth2ClientCreateFirstRequest) (*dto.OAuth2ClientCreateByAdminResponse, error)
}
