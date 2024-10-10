package resource

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
)

type User struct {
	ID           int64
	Username     string
	DisplayName  string
	AllowedScope string
}

func NewUser(ctx context.Context, user domain.User, needFilter bool) User {
	usecaseUser := User{
		ID:           user.ID,
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		AllowedScope: user.AllowedScope.String(),
	}

	if needFilter {
		Filter(ctx, &usecaseUser.AllowedScope).
			WhenRequestUserNot(user.ID).
			WhenNotContainsScope(domain.ScopeEngine.New(
				domain.Actions.Read,
				domain.Resources.User.AllowedScope,
			))
	}

	return usecaseUser
}
