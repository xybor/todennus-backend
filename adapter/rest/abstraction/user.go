package abstraction

import (
	"context"

	"github.com/xybor/todennus-backend/usecase/dto"
)

type UserUsecase interface {
	Register(ctx context.Context, req dto.UserRegisterRequestDTO) (dto.UserRegisterResponseDTO, error)
	GetByID(ctx context.Context, req dto.UserGetByIDRequestDTO) (dto.UserGetByIDResponseDTO, error)
	GetByUsername(ctx context.Context, req dto.UserGetByUsernameRequestDTO) (dto.UserGetByUsernameResponseDTO, error)
	ValidateCredentials(ctx context.Context, req dto.UserValidateCredentialsRequestDTO) (dto.UserValidateCredentialsResponseDTO, error)
}
