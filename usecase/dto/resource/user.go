package resource

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/pkg/enum"
)

type User struct {
	ID           int64
	Username     string
	DisplayName  string
	AllowedScope string
	Role         enum.Enum[domain.UserRole]
}

func NewUser(ctx context.Context, user domain.User, needFilter bool) User {
	usecaseUser := User{
		ID:           user.ID,
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		AllowedScope: user.AllowedScope.String(),
		Role:         user.Role,
	}

	if needFilter {
		Filter(ctx, &usecaseUser.AllowedScope).
			WhenRequestUserNot(user.ID).
			WhenNotContainsScope(domain.ScopeEngine.New(
				domain.Actions.Read,
				domain.Resources.User.AllowedScope,
			))

		Set(ctx, &usecaseUser.Role, enum.Default[domain.UserRole]()).
			WhenRequestUserNot(user.ID).
			WhenNotContainsScope(domain.ScopeEngine.New(
				domain.Actions.Read,
				domain.Resources.User.Role,
			))
	}

	return usecaseUser
}
