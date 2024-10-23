package resource

import (
	"github.com/xybor/todennus-backend/usecase/dto/resource"
)

type User struct {
	ID          string `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Role        string `json:"role,omitempty"`
}

func NewUser(user *resource.User) *User {
	return &User{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role.String(),
	}
}
