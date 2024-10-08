package abstraction

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, userID int64) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshTokenID, accessTokenID int64, seq int) error
	UpdateByRefreshTokenID(ctx context.Context, refreshTokenID, accessTokenId int64, expectedCurSeq int) error
	DeleteByRefreshTokenID(ctx context.Context, refreshTokenID int64) error
}

type OAuth2ClientRepository interface {
	Create(ctx context.Context, client domain.OAuth2Client) error
	GetByID(ctx context.Context, clientID int64) (domain.OAuth2Client, error)
}
