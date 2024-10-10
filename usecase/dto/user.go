package dto

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/usecase/dto/resource"
)

// Register
type UserRegisterRequestDTO struct {
	Username string
	Password string
}

type UserRegisterResponseDTO struct {
	User resource.User
}

func NewUserRegisterResponseDTO(ctx context.Context, user domain.User) UserRegisterResponseDTO {
	return UserRegisterResponseDTO{
		User: resource.NewUser(ctx, user, false),
	}
}

// GetByID
type UserGetByIDRequestDTO struct {
	UserID int64
}

type UserGetByIDResponseDTO struct {
	User resource.User
}

func NewUserGetByIDResponseDTO(ctx context.Context, user domain.User) UserGetByIDResponseDTO {
	return UserGetByIDResponseDTO{
		User: resource.NewUser(ctx, user, true),
	}
}

// GetByUsername
type UserGetByUsernameRequestDTO struct {
	Username string
}

type UserGetByUsernameResponseDTO struct {
	User resource.User
}

func NewUserGetByUsernameResponseDTO(ctx context.Context, user domain.User) UserGetByUsernameResponseDTO {
	return UserGetByUsernameResponseDTO{
		User: resource.NewUser(ctx, user, true),
	}
}
