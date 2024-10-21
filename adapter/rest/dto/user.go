package dto

import (
	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/adapter/rest/dto/resource"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/todennus-backend/usecase/dto"
	"github.com/xybor/x/xerror"
)

func ParseUserID(meID snowflake.ID, s string) (snowflake.ID, error) {
	if s == "@me" {
		return meID, nil
	}

	return snowflake.ParseString(s)
}

// Register
type UserRegisterRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req UserRegisterRequestDTO) To() *dto.UserRegisterRequestDTO {
	return &dto.UserRegisterRequestDTO{
		Username: req.Username,
		Password: req.Password,
	}
}

type UserRegisterResponseDTO struct {
	*resource.User
}

func NewUserRegisterResponseDTO(resp *dto.UserRegisterResponseDTO) *UserRegisterResponseDTO {
	return &UserRegisterResponseDTO{
		User: resource.NewUser(resp.User),
	}
}

// GetByID
type UserGetByIDRequestDTO struct {
	UserID string `param:"user_id"`
}

func (req UserGetByIDRequestDTO) To(meID snowflake.ID) (*dto.UserGetByIDRequestDTO, error) {
	userID, err := ParseUserID(meID, req.UserID)
	if err != nil {
		return nil, xerror.Enrich(usecase.ErrRequestInvalid, "user id is invalid").
			Hide(err, "failed-to-parse-user-id", "uid", req.UserID)
	}

	return &dto.UserGetByIDRequestDTO{UserID: userID}, nil
}

type UserGetByIDResponseDTO struct {
	*resource.User
}

func NewUserGetByIDResponseDTO(resp *dto.UserGetByIDResponseDTO) *UserGetByIDResponseDTO {
	return &UserGetByIDResponseDTO{
		User: resource.NewUser(resp.User),
	}
}

// GetByUsername
type UserGetByUsernameRequestDTO struct {
	Username string `param:"username"`
}

func (req UserGetByUsernameRequestDTO) To() *dto.UserGetByUsernameRequestDTO {
	return &dto.UserGetByUsernameRequestDTO{
		Username: req.Username,
	}
}

type UserGetByUsernameResponseDTO struct {
	*resource.User
}

func NewUserGetByUsernameResponseDTO(resp *dto.UserGetByUsernameResponseDTO) *UserGetByUsernameResponseDTO {
	return &UserGetByUsernameResponseDTO{
		User: resource.NewUser(resp.User),
	}
}

// Validate
type UserValidateRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req UserValidateRequestDTO) To() *dto.UserValidateCredentialsRequestDTO {
	return &dto.UserValidateCredentialsRequestDTO{
		Username: req.Username,
		Password: req.Password,
	}
}

type UserValidateResponseDTO struct {
	*resource.User
}

func NewUserValidateResponseDTO(resp *dto.UserValidateCredentialsResponseDTO) *UserValidateResponseDTO {
	return &UserValidateResponseDTO{
		User: resource.NewUser(resp.User),
	}
}
