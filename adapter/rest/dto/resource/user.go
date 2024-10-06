package resource

import (
	"strconv"

	"github.com/xybor/todennus-backend/domain"
)

type User struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

func NewUser(user domain.User) User {
	return User{
		ID:          strconv.FormatInt(user.ID, 10),
		Username:    user.Username,
		DisplayName: user.DisplayName,
	}
}
