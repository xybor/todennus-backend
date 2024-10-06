package abstraction

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
)

type UserRepository interface {
	// Save inserts user if it doesn't exist, otherwise update.
	Save(ctx context.Context, user domain.User) error

	// GetByUsername returns the user by username.
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}
