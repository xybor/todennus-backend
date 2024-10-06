package abstraction

import (
	"context"

	"github.com/xybor/todennus-backend/usecase/dto"
)

type UserUsecase interface {
	Register(ctx context.Context, req dto.UserRegisterRequestDTO) (dto.UserRegisterResponseDTO, error)
	Validate(ctx context.Context, req dto.UserValidateRequestDTO) (dto.UserValidateResponseDTO, error)
}

type OAuth2Usecase interface {
	Token(ctx context.Context, req dto.OAuth2TokenRequest) (dto.OAuth2TokenResponse, error)
}
