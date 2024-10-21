package resource

import (
	"context"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/x/enum"
	"github.com/xybor/x/scope"
)

type User struct {
	ID           snowflake.ID
	Username     string
	DisplayName  string
	AllowedScope string
	Role         enum.Enum[domain.UserRole]
}

func NewUser(ctx context.Context, user *domain.User, needFilter bool) *User {
	usecaseUser := &User{
		ID:           user.ID,
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		AllowedScope: user.AllowedScope.String(),
		Role:         user.Role,
	}

	if needFilter {
		Filter(ctx, &usecaseUser.AllowedScope).
			WhenRequestUserNot(user.ID).
			WhenNotContainsScope(scope.New(domain.Actions.Read, domain.Resources.User.AllowedScope))

		Set(ctx, &usecaseUser.Role, enum.Default[domain.UserRole]()).
			WhenRequestUserNot(user.ID).
			WhenNotContainsScope(scope.New(domain.Actions.Read, domain.Resources.User.Role))
	}

	return usecaseUser
}
