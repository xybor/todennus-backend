package resource

import (
	"github.com/xybor/todennus-backend/usecase/dto/resource"
)

type User struct {
	ID          string `json:"id,omitempty" example:"330559330522759168"`
	Username    string `json:"username,omitempty" example:"huykingsofm"`
	DisplayName string `json:"display_name,omitempty" example:"Huy Le Ngoc"`
	Role        string `json:"role,omitempty" example:"admin"`
}

func NewUser(user *resource.User) *User {
	return &User{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role.String(),
	}
}
