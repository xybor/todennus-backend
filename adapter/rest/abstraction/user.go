package abstraction

import (
	"context"

	"github.com/xybor/todennus-backend/usecase/dto"
)

type UserUsecase interface {
	Register(ctx context.Context, req dto.UserRegisterRequestDTO) (dto.UserRegisterResponseDTO, error)
	GetByID(ctx context.Context, req dto.UserGetByIDRequestDTO) (dto.UserGetByIDResponseDTO, error)
	GetByUsername(ctx context.Context, req dto.UserGetByUsernameRequestDTO) (dto.UserGetByUsernameResponseDTO, error)
}

type OAuth2Usecase interface {
	Token(ctx context.Context, req dto.OAuth2TokenRequestDTO) (dto.OAuth2TokenResponseDTO, error)
}

type OAuth2ClientUsecase interface {
	CreateClient(ctx context.Context, req dto.OAuth2ClientCreateRequestDTO) (dto.OAuth2ClientCreateResponseDTO, error)
	GetClient(ctx context.Context, req dto.OAuth2ClientGetRequestDTO) (dto.OAuth2ClientGetResponseDTO, error)
}
