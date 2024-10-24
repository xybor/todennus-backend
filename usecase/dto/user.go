package dto

import (
	"context"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/usecase/dto/resource"
)

// Register
type UserRegisterRequest struct {
	Username string
	Password string
}

type UserRegisterResponse struct {
	User *resource.User
}

func NewUserRegisterResponse(ctx context.Context, user *domain.User) *UserRegisterResponse {
	return &UserRegisterResponse{
		User: resource.NewUserWithoutFilter(user),
	}
}

// GetByID
type UserGetByIDRequest struct {
	UserID snowflake.ID
}

type UserGetByIDResponse struct {
	User *resource.User
}

func NewUserGetByIDResponse(ctx context.Context, user *domain.User) *UserGetByIDResponse {
	return &UserGetByIDResponse{
		User: resource.NewUser(ctx, user),
	}
}

// GetByUsername
type UserGetByUsernameRequest struct {
	Username string
}

type UserGetByUsernameResponse struct {
	User *resource.User
}

func NewUserGetByUsernameResponse(ctx context.Context, user *domain.User) *UserGetByUsernameResponse {
	return &UserGetByUsernameResponse{
		User: resource.NewUser(ctx, user),
	}
}

// Validate
type UserValidateCredentialsRequest struct {
	Username string
	Password string
}

type UserValidateCredentialsResponse struct {
	User *resource.User
}

func NewUserValidateCredentialsResponse(ctx context.Context, user *domain.User) *UserValidateCredentialsResponse {
	return &UserValidateCredentialsResponse{
		User: resource.NewUser(ctx, user),
	}
}
