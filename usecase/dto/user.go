package dto

import "github.com/xybor/todennus-backend/domain"

type UserRegisterRequestDTO struct {
	Username string
	Password string
}

type UserRegisterResponseDTO struct {
	User domain.User
}

type UserRegisterSuperAdminRequestDTO struct {
	Username string
	Password string
}

type UserRegisterSuperAdminResponseDTO struct {
	User domain.User
}
