package dto

import (
	"context"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/adapter/rest/dto/resource"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/todennus-backend/usecase/dto"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xerror"
)

func ParseUserID(ctx context.Context, s string) (snowflake.ID, error) {
	if s == "@me" {
		return xcontext.RequestUserID(ctx), nil
	}

	return snowflake.ParseString(s)
}

// Register
type UserRegisterRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req UserRegisterRequestDTO) To() dto.UserRegisterRequestDTO {
	return dto.UserRegisterRequestDTO{
		Username: req.Username,
		Password: req.Password,
	}
}

type UserRegisterResponseDTO struct {
	resource.User
}

func NewUserRegisterResponseDTO(resp dto.UserRegisterResponseDTO) UserRegisterResponseDTO {
	return UserRegisterResponseDTO{
		User: resource.NewUser(resp.User),
	}
}

// GetByID
type UserGetByIDRequestDTO struct {
	UserID string `param:"user_id"`
}

func (req UserGetByIDRequestDTO) To(ctx context.Context) (dto.UserGetByIDRequestDTO, error) {
	userID, err := ParseUserID(ctx, req.UserID)
	if err != nil {
		xcontext.Logger(ctx).Debug("failed-to-parse-user-id", "err", err, "uid", req.UserID)
		return dto.UserGetByIDRequestDTO{}, xerror.Wrap(usecase.ErrRequestInvalid, "user id is invalid")
	}

	return dto.UserGetByIDRequestDTO{UserID: userID}, nil
}

type UserGetByIDResponseDTO struct {
	resource.User
}

func NewUserGetByIDResponseDTO(resp dto.UserGetByIDResponseDTO) UserGetByIDResponseDTO {
	return UserGetByIDResponseDTO{
		User: resource.NewUser(resp.User),
	}
}

// GetByUsername
type UserGetByUsernameRequestDTO struct {
	Username string `param:"username"`
}

func (req UserGetByUsernameRequestDTO) To() dto.UserGetByUsernameRequestDTO {
	return dto.UserGetByUsernameRequestDTO{
		Username: req.Username,
	}
}

type UserGetByUsernameResponseDTO struct {
	resource.User
}

func NewUserGetByUsernameResponseDTO(resp dto.UserGetByUsernameResponseDTO) UserGetByUsernameResponseDTO {
	return UserGetByUsernameResponseDTO{
		User: resource.NewUser(resp.User),
	}
}
