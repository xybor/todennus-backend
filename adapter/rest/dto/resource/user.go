package resource

import (
	"github.com/xybor/todennus-backend/pkg/xstring"
	"github.com/xybor/todennus-backend/usecase/dto/resource"
)

type User struct {
	ID           string `json:"id,omitempty"`
	Username     string `json:"username,omitempty"`
	DisplayName  string `json:"display_name,omitempty"`
	AllowedScope string `json:"allowed_scope,omitempty"`
}

func NewUser(user resource.User) User {
	return User{
		ID:           xstring.FormatID(user.ID),
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		AllowedScope: user.AllowedScope,
	}
}
