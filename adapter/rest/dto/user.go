package dto

import (
	"github.com/xybor/todennus-backend/adapter/rest/dto/resource"
	"github.com/xybor/todennus-backend/usecase/dto"
)

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

func NewUserRegisterResponse(resp dto.UserRegisterResponseDTO) UserRegisterResponseDTO {
	return UserRegisterResponseDTO{
		User: resource.NewUser(resp.User),
	}
}
