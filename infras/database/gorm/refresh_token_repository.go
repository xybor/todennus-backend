package gorm

import (
	"context"

	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/infras/database/model"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (repo *RefreshTokenRepository) Create(
	ctx context.Context,
	refreshTokenId,
	accessTokenID int64,
	seq int,
) error {
	return database.ConvertError(repo.db.WithContext(ctx).Create(&model.RefreshTokenModel{
		RefreshTokenID: refreshTokenId,
		AccessTokenID:  accessTokenID,
		Seq:            seq,
	}).Error)
}

func (repo *RefreshTokenRepository) UpdateByRefreshTokenID(
	ctx context.Context,
	refreshTokenID, accessTokenID int64,
	expectedCurSeq int,
) error {
	result := repo.db.WithContext(ctx).Model(&model.RefreshTokenModel{}).
		Where("refresh_token_id=? AND seq=?", refreshTokenID, expectedCurSeq).
		Updates(map[string]any{
			"seq":             expectedCurSeq + 1,
			"access_token_id": accessTokenID,
		})

	if result.RowsAffected == 0 {
		return database.ErrRecordNotFound
	}

	return database.ConvertError(result.Error)
}

func (repo *RefreshTokenRepository) DeleteByRefreshTokenID(
	ctx context.Context, refreshTokenID int64,
) error {
	return database.ConvertError(repo.db.WithContext(ctx).Delete(&model.RefreshTokenModel{}, refreshTokenID).Error)
}
