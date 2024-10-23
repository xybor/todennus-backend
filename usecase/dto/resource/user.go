package resource

import (
	"context"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/x/enum"
	"github.com/xybor/x/scope"
)

type User struct {
	ID          snowflake.ID
	Username    string
	DisplayName string
	Role        enum.Enum[domain.UserRole]
}

func NewUser(ctx context.Context, user *domain.User) *User {
	usecaseUser := &User{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role,
	}

	Set(ctx, &usecaseUser.Role, enum.Default[domain.UserRole]()).
		WhenRequestUserNot(user.ID).
		WhenNotContainsScope(scope.New(domain.Actions.Read, domain.Resources.User.Role))

	return usecaseUser
}

func NewUserWithoutFilter(user *domain.User) *User {
	usecaseUser := &User{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role,
	}

	return usecaseUser
}
