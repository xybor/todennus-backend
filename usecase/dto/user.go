package dto

import "github.com/xybor/todennus-backend/domain"

type UserRegisterRequestDTO struct {
	Username string
	Password string
}

type UserRegisterResponseDTO struct {
	User domain.User
}

type UserValidateRequestDTO struct {
	Username string
	Password string
}

type UserValidateResponseDTO struct {
	User domain.User
}
